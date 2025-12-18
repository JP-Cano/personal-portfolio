import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { Project, Projects } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<Projects>(() =>
    api.get<Projects>(API_PATHS.PROJECTS, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch projects";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }

  return new Response(JSON.stringify({ data }), {
    status: 200,
    headers: { "Content-Type": "application/json" },
  });
};

export const POST: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";
  const body = await request.json();

  const [data, error] = await asyncThrowable<Project>(() =>
    api.post<Project>(API_PATHS.PROJECTS, body, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to create project";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }

  return new Response(JSON.stringify({ data }), {
    status: 201,
    headers: { "Content-Type": "application/json" },
  });
};
