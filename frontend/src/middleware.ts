import { defineMiddleware } from "astro:middleware";
import { api } from "@/api/api-client";
import type { User } from "@/types/auth.ts";
import { asyncThrowable } from "@/utils/utils.ts";

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

  const [user, error] = await asyncThrowable<User>(() =>
    api.getCurrentUser(cookie)
  );

  if (error) {
    return context.redirect("/admin/login");
  }

  context.locals.user = user as User;

  return next();
});
