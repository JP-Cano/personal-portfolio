import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable(() =>
    api.get(`${API_PATHS.UPLOAD_CERTIFICATES}/${id}`, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch certification";
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
    api.delete(`${API_PATHS.UPLOAD_CERTIFICATES}/${id}`, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to delete certification";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }

  return new Response(
    JSON.stringify({ data: { message: "Certification deleted" } }),
    {
      status: 200,
      headers: { "Content-Type": "application/json" },
    }
  );
};
