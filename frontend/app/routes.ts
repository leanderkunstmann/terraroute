import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("globe", "routes/globe.tsx"),
  route("*", "routes/redirect.tsx"),
] satisfies RouteConfig;
