import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { ExperienceClient, ExperienceClients } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<ExperienceClients>(() =>
    api.get<ExperienceClients>(API_PATHS.EXPERIENCE_CLIENTS(id ?? ""), cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch clients";
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

export const POST: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";
  const body = await request.json();

  const [data, error] = await asyncThrowable<ExperienceClient>(() =>
    api.post<ExperienceClient>(
      API_PATHS.EXPERIENCE_CLIENTS(id ?? ""),
      body,
      cookie
    )
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to create client";
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
