import { useState, useRef, useEffect, useCallback, useMemo } from "react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { useGetAssetsQuery } from "@/modules/asset/actions/assetApiNew";
import { cn } from "@/lib/utils";
import { Search } from "lucide-react";

type AssetSelectProps = {
  value: string;
  onChange: (value: string) => void;
  excludeIds?: string[];
  disabled?: boolean;
};

export function AssetSelect({ value, onChange, excludeIds = [], disabled }: AssetSelectProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [activeIndex, setActiveIndex] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  const excludeSet = useMemo(() => new Set(excludeIds), [excludeIds]);

  const { data, isFetching } = useGetAssetsQuery(
    { search: search || undefined, limit: 20 },
    { skip: !open }
  );

  const allAssets = data?.data ?? [];
  const assets = useMemo(
    () => allAssets.filter((a) => !excludeSet.has(a.id)),
    [allAssets, excludeSet]
  );
  const selectedAsset = allAssets.find((a) => a.id === value) ?? allAssets[0];

  const handleSelect = useCallback(
    (assetId: string) => {
      onChange(assetId);
      setSearch("");
      setOpen(false);
      setActiveIndex(-1);
    },
    [onChange]
  );

  useEffect(() => {
    if (!open) {
      setSearch("");
      setActiveIndex(-1);
    }
  }, [open]);

  useEffect(() => {
    setActiveIndex(-1);
  }, [search]);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!open) return;

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        setActiveIndex((prev) => (prev < assets.length - 1 ? prev + 1 : 0));
        break;
      case "ArrowUp":
        e.preventDefault();
        setActiveIndex((prev) => (prev > 0 ? prev - 1 : assets.length - 1));
        break;
      case "Enter":
        e.preventDefault();
        if (activeIndex >= 0 && assets[activeIndex]) {
          handleSelect(assets[activeIndex].id);
        }
        break;
      case "Escape":
        setOpen(false);
        break;
    }
  };

  useEffect(() => {
    if (activeIndex >= 0 && listRef.current) {
      const item = listRef.current.children[activeIndex] as HTMLElement;
      item?.scrollIntoView({ block: "nearest" });
    }
  }, [activeIndex]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          disabled={disabled}
          className={cn(
            "flex h-9 w-full items-center justify-between rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
            !selectedAsset && "text-muted-foreground"
          )}
          onClick={() => setOpen(true)}
        >
          {selectedAsset ? (
            <span>
              {selectedAsset.code} — {selectedAsset.name}
            </span>
          ) : (
            <span>Pilih aset...</span>
          )}
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-[350px] p-0" align="start">
        <div className="flex items-center border-b px-3">
          <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
          <Input
            ref={inputRef}
            placeholder="Cari kode atau nama aset..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onKeyDown={handleKeyDown}
            className="h-9 border-0 shadow-none focus-visible:ring-0"
            autoFocus
          />
        </div>
        <div ref={listRef} className="max-h-[300px] overflow-y-auto p-1">
          {isFetching && !assets.length ? (
            <div className="py-6 text-center text-sm text-muted-foreground">
              Loading...
            </div>
          ) : assets.length === 0 ? (
            <div className="py-6 text-center text-sm text-muted-foreground">
              {search ? "Tidak ada aset ditemukan" : "Ketik untuk mencari aset"}
            </div>
          ) : (
            assets.map((asset, index) => (
              <button
                key={asset.id}
                type="button"
                className={cn(
                  "flex w-full items-start gap-2 rounded-sm px-2 py-1.5 text-left text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                  activeIndex === index && "bg-accent text-accent-foreground"
                )}
                onMouseDown={(e) => {
                  e.preventDefault();
                  handleSelect(asset.id);
                }}
                onMouseEnter={() => setActiveIndex(index)}
              >
                <div className="flex-1 min-w-0">
                  <div className="font-medium truncate">{asset.code}</div>
                  <div className="text-muted-foreground truncate">{asset.name}</div>
                </div>
                {asset.category && (
                  <span className="text-xs text-muted-foreground whitespace-nowrap">
                    {asset.category.name}
                  </span>
                )}
              </button>
            ))
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
