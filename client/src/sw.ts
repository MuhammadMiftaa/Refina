import { BackgroundSyncPlugin } from "workbox-background-sync";
import { cleanupOutdatedCaches, precacheAndRoute } from "workbox-precaching";
import { NavigationRoute, Route, registerRoute } from "workbox-routing";
import { CacheFirst, NetworkFirst, NetworkOnly } from "workbox-strategies";
import { getViteAPIBaseURL } from "./lib/readenv";

declare let self: ServiceWorkerGlobalScope;

cleanupOutdatedCaches();

precacheAndRoute(self.__WB_MANIFEST);

// self.skipWaiting();

// cache images
const imageRoute = new Route(
  ({ request, sameOrigin }) => {
    return sameOrigin && request.destination === "image";
  },
  new CacheFirst({
    cacheName: "images",
  }),
);
registerRoute(imageRoute);

// cache fonts
const fontRoute = new Route(
  ({ request, sameOrigin }) => {
    return sameOrigin && request.destination === "font";
  },
  new CacheFirst({
    cacheName: "fonts",
  }),
);
registerRoute(fontRoute);

// cache api calls
const fetchTasksRoute = new Route(
  ({ request }) => {
    return request.url === getViteAPIBaseURL() + "/tasks";
  },
  new NetworkFirst({
    cacheName: "api/fetch-tasks",
  }),
);
registerRoute(fetchTasksRoute);

// cache navigations
const navigationRoute = new NavigationRoute(
  new NetworkFirst({
    cacheName: "navigation",
    networkTimeoutSeconds: 3,
  }),
);
registerRoute(navigationRoute);

//  background sync
const bgSyncPlugin = new BackgroundSyncPlugin("backgroundSyncQueue", {
  maxRetentionTime: 24 * 60,
});

const taskSubmitRoute = new Route(
  ({ request }) => {
    return request.url === getViteAPIBaseURL() + "/task/create";
  },
  new NetworkOnly({
    plugins: [bgSyncPlugin],
  }),
  "POST",
);
registerRoute(taskSubmitRoute);

const editTaskRoute = new Route(
  ({ request }) => {
    return request.url.includes(getViteAPIBaseURL() + "/task");
  },
  new NetworkOnly({
    plugins: [bgSyncPlugin],
  }),
  "PATCH",
);
registerRoute(editTaskRoute);
