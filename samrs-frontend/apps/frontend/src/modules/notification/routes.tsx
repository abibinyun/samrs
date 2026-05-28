import { tenantRoute } from "@/routes";
import { createRoute } from "@tanstack/react-router";
import { NotificationPage } from "./NotificationPage";
import { Can } from "../auth/components/Can";
import { SCOPES } from "@/constants";

export const notificationRoutes = [
  createRoute({ 
    getParentRoute: () => tenantRoute, 
    path: "/notifications", 
    component: () => <Can perform={SCOPES.NOTIFICATION.SEND}><NotificationPage /></Can>
  }),
  // createRoute({ getParentRoute: () => tenantRoute, path: "/assets/$assetId", component: AssetDetailPage }),
];