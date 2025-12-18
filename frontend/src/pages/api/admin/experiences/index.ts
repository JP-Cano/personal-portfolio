import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { Experience, Experiences } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<Experiences>(() =>
    api.get<Experiences>(API_PATHS.EXPERIENCES, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch experiences";
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

  const [data, error] = await asyncThrowable<Experience>(() =>
    api.post<Experience>(API_PATHS.EXPERIENCES, body, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to create experience";
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
