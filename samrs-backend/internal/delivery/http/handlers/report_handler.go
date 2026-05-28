package handlers

import (
	"samrs-backend/internal/delivery/http/httputil"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"samrs-backend/internal/domain"
	"samrs-backend/internal/repository"
	reportusecase "samrs-backend/internal/usecase/report"
	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportHandler struct {
	usecase reportusecase.ReportUsecase
}

func NewReportHandler(u reportusecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{usecase: u}
}

func (h *ReportHandler) ExportAssets(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	createdRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}
	purchaseRange, err := httputil.ParseDateRange(c, "purchase_from", "purchase_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format purchase_date tidak valid", err)
		return
	}

	filter := repository.AssetFilter{
		Status:       strings.TrimSpace(c.Query("status")),
		Search:       strings.TrimSpace(c.Query("search")),
		DateFrom:     createdRange.From,
		DateTo:       createdRange.To,
		PurchaseFrom: purchaseRange.From,
		PurchaseTo:   purchaseRange.To,
		SortBy:       strings.TrimSpace(c.Query("sort_by")),
		SortDir:      strings.TrimSpace(c.Query("sort_dir")),
	}

	if categoryStr := strings.TrimSpace(c.Query("category_id")); categoryStr != "" {
		id, err := strconv.ParseUint(categoryStr, 10, 32)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format category_id tidak valid", err)
			return
		}
		filter.CategoryID = uint(id)
	}
	if roomStr := strings.TrimSpace(c.Query("room_id")); roomStr != "" {
		roomID, err := uuid.Parse(roomStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format room_id tidak valid", err)
			return
		}
		filter.RoomID = roomID
	}
	if bedStr := strings.TrimSpace(c.Query("bed_id")); bedStr != "" {
		bedID, err := strconv.ParseUint(bedStr, 10, 32)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format bed_id tidak valid", err)
			return
		}
		filter.BedID = uint(bedID)
	}

	items, err := h.usecase.ExportAssets(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal export assets", err)
		return
	}

	fields, headers, err := resolveAssetExportFields(c)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format export assets tidak valid", err)
		return
	}

	rows := buildAssetRows(fields, items)

	if err := writeReportExport(c, "assets", headers, rows); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengirim export assets", err)
		return
	}
}

func (h *ReportHandler) ExportComplaints(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	createdRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}

	filter := repository.ComplaintFilter{
		Status:   strings.TrimSpace(c.Query("status")),
		Search:   strings.TrimSpace(c.Query("search")),
		DateFrom: createdRange.From,
		DateTo:   createdRange.To,
		SortBy:   strings.TrimSpace(c.Query("sort_by")),
		SortDir:  strings.TrimSpace(c.Query("sort_dir")),
	}

	if assetStr := strings.TrimSpace(c.Query("asset_id")); assetStr != "" {
		assetID, err := uuid.Parse(assetStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format asset_id tidak valid", err)
			return
		}
		filter.AssetID = assetID
	}
	if assignedStr := strings.TrimSpace(c.Query("assigned_to")); assignedStr != "" {
		assignedID, err := uuid.Parse(assignedStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format assigned_to tidak valid", err)
			return
		}
		filter.AssignedTo = assignedID
	}
	if reportedStr := strings.TrimSpace(c.Query("reported_by")); reportedStr != "" {
		reportedID, err := uuid.Parse(reportedStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format reported_by tidak valid", err)
			return
		}
		filter.ReportedBy = reportedID
	}

	items, err := h.usecase.ExportComplaints(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal export complaints", err)
		return
	}

	fields, headers, err := resolveComplaintExportFields(c)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format export complaints tidak valid", err)
		return
	}

	rows := buildComplaintRows(fields, items)

	if err := writeReportExport(c, "complaints", headers, rows); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengirim export complaints", err)
		return
	}
}

func (h *ReportHandler) ExportMaintenance(c *gin.Context) {
	tenantID, ok := httputil.TenantIDFromContext(c)
	if !ok {
		util.ErrorResponse(c, http.StatusUnauthorized, "Tenant ID tidak ditemukan", nil)
		return
	}

	createdRange, err := httputil.ParseDateRange(c, "date_from", "date_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format tanggal tidak valid", err)
		return
	}
	dueRange, err := httputil.ParseDateRange(c, "due_from", "due_to")
	if err != nil {
		util.ErrorResponseFromErr(c, "Format due date tidak valid", err)
		return
	}

	filter := repository.MaintenanceScheduleFilter{
		ScheduleType: strings.TrimSpace(c.Query("schedule_type")),
		Status:       strings.TrimSpace(c.Query("status")),
		DateFrom:     createdRange.From,
		DateTo:       createdRange.To,
		DueFrom:      dueRange.From,
		DueTo:        dueRange.To,
		SortBy:       strings.TrimSpace(c.Query("sort_by")),
		SortDir:      strings.TrimSpace(c.Query("sort_dir")),
	}

	if assetStr := strings.TrimSpace(c.Query("asset_id")); assetStr != "" {
		assetID, err := uuid.Parse(assetStr)
		if err != nil {
			util.ErrorResponseFromErr(c, "Format asset_id tidak valid", err)
			return
		}
		filter.AssetID = assetID
	}

	items, err := h.usecase.ExportMaintenance(tenantID, filter)
	if err != nil {
		util.ErrorResponseFromErr(c, "Gagal export schedules", err)
		return
	}

	fields, headers, err := resolveMaintenanceExportFields(c)
	if err != nil {
		util.ErrorResponseFromErr(c, "Format export schedules tidak valid", err)
		return
	}

	rows := buildMaintenanceRows(fields, items)

	if err := writeReportExport(c, "maintenance_schedules", headers, rows); err != nil {
		util.ErrorResponseFromErr(c, "Gagal mengirim export schedules", err)
		return
	}
}

func writeCSVHeaders(c *gin.Context, filename string) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename="+filename)
}

func writeReportExport(c *gin.Context, prefix string, headers []string, rows [][]string) error {
	format := strings.ToLower(strings.TrimSpace(c.Query("format")))
	if format == "" {
		format = "csv"
	}

	switch format {
	case "csv":
		filename := fmt.Sprintf("%s_%s.csv", prefix, time.Now().Format("20060102_150405"))
		writeCSVHeaders(c, filename)
		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()
		if err := writer.Write(headers); err != nil {
			return err
		}
		for _, row := range rows {
			if err := writer.Write(row); err != nil {
				return err
			}
		}
		return nil
	case "xlsx":
		filename := fmt.Sprintf("%s_%s.xlsx", prefix, time.Now().Format("20060102_150405"))
		buf, err := buildXLSX(headers, rows)
		if err != nil {
			return err
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		_, err = c.Writer.Write(buf.Bytes())
		return err
	case "pdf":
		filename := fmt.Sprintf("%s_%s.pdf", prefix, time.Now().Format("20060102_150405"))
		pdfBytes, err := buildPDF(prefix, headers, rows)
		if err != nil {
			return err
		}
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename="+filename)
		_, err = c.Writer.Write(pdfBytes)
		return err
	default:
		return errors.New("format tidak didukung (gunakan csv, xlsx, atau pdf)")
	}
}

func buildXLSX(headers []string, rows [][]string) (*bytes.Buffer, error) {
	file := excelize.NewFile()
	sheet := file.GetSheetName(0)
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err := file.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}
	for rowIdx, row := range rows {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			if err := file.SetCellValue(sheet, cell, value); err != nil {
				return nil, err
			}
		}
	}
	buf, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func buildPDF(title string, headers []string, rows [][]string) ([]byte, error) {
	orientation := "P"
	if len(headers) > 8 {
		orientation = "L"
	}
	pdf := gofpdf.New(orientation, "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Report %s", strings.ReplaceAll(title, "_", " ")))
	pdf.Ln(10)

	pageWidth, _ := pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	tableWidth := pageWidth - leftMargin - rightMargin
	colWidth := tableWidth / float64(len(headers))

	pdf.SetFont("Arial", "B", 8)
	for _, header := range headers {
		pdf.CellFormat(colWidth, 6, header, "1", 0, "LM", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, row := range rows {
		for _, value := range row {
			pdf.CellFormat(colWidth, 6, value, "1", 0, "LM", false, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02")
}

type exportField struct {
	Key   string
	Label string
}

func resolveAssetExportFields(c *gin.Context) ([]exportField, []string, error) {
	defaults := []exportField{
		{Key: "id", Label: "id"},
		{Key: "code", Label: "code"},
		{Key: "name", Label: "name"},
		{Key: "status", Label: "status"},
		{Key: "category", Label: "category"},
		{Key: "room", Label: "room"},
		{Key: "bed", Label: "bed"},
		{Key: "vendor", Label: "vendor"},
		{Key: "brand", Label: "brand"},
		{Key: "model", Label: "model"},
		{Key: "purchase_date", Label: "purchase_date"},
		{Key: "created_at", Label: "created_at"},
	}
	allowed := make(map[string]exportField, len(defaults))
	for _, field := range defaults {
		allowed[field.Key] = field
	}
	fields, err := parseExportFields(c.Query("fields"), defaults, allowed)
	if err != nil {
		return nil, nil, err
	}
	headers := applyLabelOverrides(fields, parseColumnLabels(c.Query("columns")))
	return fields, headers, nil
}

func resolveComplaintExportFields(c *gin.Context) ([]exportField, []string, error) {
	defaults := []exportField{
		{Key: "id", Label: "id"},
		{Key: "asset_code", Label: "asset_code"},
		{Key: "title", Label: "title"},
		{Key: "status", Label: "status"},
		{Key: "reported_by", Label: "reported_by"},
		{Key: "assigned_to", Label: "assigned_to"},
		{Key: "created_at", Label: "created_at"},
	}
	allowed := make(map[string]exportField, len(defaults))
	for _, field := range defaults {
		allowed[field.Key] = field
	}
	fields, err := parseExportFields(c.Query("fields"), defaults, allowed)
	if err != nil {
		return nil, nil, err
	}
	headers := applyLabelOverrides(fields, parseColumnLabels(c.Query("columns")))
	return fields, headers, nil
}

func resolveMaintenanceExportFields(c *gin.Context) ([]exportField, []string, error) {
	defaults := []exportField{
		{Key: "id", Label: "id"},
		{Key: "asset_code", Label: "asset_code"},
		{Key: "schedule_type", Label: "schedule_type"},
		{Key: "title", Label: "title"},
		{Key: "interval_days", Label: "interval_days"},
		{Key: "next_due_date", Label: "next_due_date"},
		{Key: "last_done_at", Label: "last_done_at"},
		{Key: "status", Label: "status"},
		{Key: "created_at", Label: "created_at"},
	}
	allowed := make(map[string]exportField, len(defaults))
	for _, field := range defaults {
		allowed[field.Key] = field
	}
	fields, err := parseExportFields(c.Query("fields"), defaults, allowed)
	if err != nil {
		return nil, nil, err
	}
	headers := applyLabelOverrides(fields, parseColumnLabels(c.Query("columns")))
	return fields, headers, nil
}

func parseExportFields(param string, defaults []exportField, allowed map[string]exportField) ([]exportField, error) {
	if strings.TrimSpace(param) == "" {
		return defaults, nil
	}
	parts := strings.Split(param, ",")
	fields := make([]exportField, 0, len(parts))
	for _, part := range parts {
		key := strings.TrimSpace(part)
		if key == "" {
			continue
		}
		field, ok := allowed[key]
		if !ok {
			return nil, fmt.Errorf("field tidak dikenal: %s", key)
		}
		fields = append(fields, field)
	}
	if len(fields) == 0 {
		return defaults, nil
	}
	return fields, nil
}

func parseColumnLabels(param string) map[string]string {
	overrides := map[string]string{}
	if strings.TrimSpace(param) == "" {
		return overrides
	}
	pairs := strings.Split(param, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		label := strings.TrimSpace(parts[1])
		if key == "" || label == "" {
			continue
		}
		overrides[key] = label
	}
	return overrides
}

func applyLabelOverrides(fields []exportField, overrides map[string]string) []string {
	headers := make([]string, 0, len(fields))
	for _, field := range fields {
		label := field.Label
		if override, ok := overrides[field.Key]; ok {
			label = override
		}
		headers = append(headers, label)
	}
	return headers
}

func buildAssetRows(fields []exportField, items []domain.Asset) [][]string {
	rows := make([][]string, 0, len(items))
	for _, asset := range items {
		row := make([]string, 0, len(fields))
		for _, field := range fields {
			row = append(row, assetFieldValue(field.Key, asset))
		}
		rows = append(rows, row)
	}
	return rows
}

func assetFieldValue(key string, asset domain.Asset) string {
	switch key {
	case "id":
		return asset.ID.String()
	case "code":
		return asset.Code
	case "name":
		return asset.Name
	case "status":
		return asset.Status
	case "category":
		return asset.Category.Name
	case "room":
		return asset.Room.Code
	case "bed":
		return asset.Bed.Code
	case "vendor":
		return asset.Vendor.Name
	case "brand":
		return asset.BrandRef.Name
	case "model":
		return asset.ModelRef.Name
	case "purchase_date":
		return formatDate(asset.PurchaseDate)
	case "created_at":
		return asset.CreatedAt.Format(time.RFC3339)
	default:
		return ""
	}
}

func buildComplaintRows(fields []exportField, items []domain.Complaint) [][]string {
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		row := make([]string, 0, len(fields))
		for _, field := range fields {
			row = append(row, complaintFieldValue(field.Key, item))
		}
		rows = append(rows, row)
	}
	return rows
}

func complaintFieldValue(key string, item domain.Complaint) string {
	switch key {
	case "id":
		return item.ID.String()
	case "asset_code":
		return item.Asset.Code
	case "title":
		return item.Title
	case "status":
		return item.Status
	case "reported_by":
		return item.Reporter.Username
	case "assigned_to":
		return item.Assignee.Username
	case "created_at":
		return item.CreatedAt.Format(time.RFC3339)
	default:
		return ""
	}
}

func buildMaintenanceRows(fields []exportField, items []domain.MaintenanceSchedule) [][]string {
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		row := make([]string, 0, len(fields))
		for _, field := range fields {
			row = append(row, maintenanceFieldValue(field.Key, item))
		}
		rows = append(rows, row)
	}
	return rows
}

func maintenanceFieldValue(key string, item domain.MaintenanceSchedule) string {
	switch key {
	case "id":
		return strconv.FormatUint(uint64(item.ID), 10)
	case "asset_code":
		return item.Asset.Code
	case "schedule_type":
		return item.ScheduleType
	case "title":
		return item.Title
	case "interval_days":
		return strconv.Itoa(item.IntervalDays)
	case "next_due_date":
		return formatDate(item.NextDueDate)
	case "last_done_at":
		return formatDate(item.LastDoneAt)
	case "status":
		return item.Status
	case "created_at":
		return item.CreatedAt.Format(time.RFC3339)
	default:
		return ""
	}
}
