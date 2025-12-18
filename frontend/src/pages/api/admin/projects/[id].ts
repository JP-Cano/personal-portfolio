import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { Project } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<Project>(() =>
    api.get<Project>(`${API_PATHS.PROJECTS}/${id}`, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch project";
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

export const PATCH: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";
  const body = await request.json();

  const [data, error] = await asyncThrowable<Project>(() =>
    api.patch<Project>(`${API_PATHS.PROJECTS}/${id}`, body, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to update project";
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

export const DELETE: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  const [_, error] = await asyncThrowable(() =>
    api.delete(`${API_PATHS.PROJECTS}/${id}`, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to delete project";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }

  return new Response(
    JSON.stringify({ data: { message: "Project deleted" } }),
    {
      status: 200,
      headers: { "Content-Type": "application/json" },
    }
  );
};
