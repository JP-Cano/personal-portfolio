import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { Experience } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<Experience>(() =>
    api.get<Experience>(`${API_PATHS.EXPERIENCES}/${id}`, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch experience";
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

  const [data, error] = await asyncThrowable<Experience>(() =>
    api.patch<Experience>(`${API_PATHS.EXPERIENCES}/${id}`, body, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to update experience";
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
    api.delete(`${API_PATHS.EXPERIENCES}/${id}`, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to delete experience";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }

  return new Response(
    JSON.stringify({ data: { message: "Experience deleted" } }),
    {
      status: 200,
      headers: { "Content-Type": "application/json" },
    }
  );
};
