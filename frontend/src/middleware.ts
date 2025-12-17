import { defineMiddleware } from "astro:middleware";
import { api } from "@/api/api-client";

const PUBLIC_ADMIN_ROUTES = new Set(["/admin/login"]);

export const onRequest = defineMiddleware(async (context, next) => {
  const { pathname } = context.url;

  if (!pathname.startsWith("/admin")) {
    return next();
  }

  const cookie = context.request.headers.get("cookie") || "";

  if (PUBLIC_ADMIN_ROUTES.has(pathname)) {
    try {
      await api.getCurrentUser(cookie);
      return context.redirect("/admin");
    } catch {
      return next();
    }
  }

  try {
    context.locals.user = await api.getCurrentUser(cookie);
    return next();
  } catch {
    return context.redirect("/admin/login");
  }
});
