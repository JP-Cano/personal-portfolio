import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { ExperienceClient } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

const clientPath = (id: string, clientId: string): string =>
  `${API_PATHS.EXPERIENCE_CLIENTS(id)}/${clientId}`;

export const GET: APIRoute = async ({ params, request }) => {
  const { id, clientId } = params;
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<ExperienceClient>(() =>
    api.get<ExperienceClient>(clientPath(id ?? "", clientId ?? ""), cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch client";
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
  const { id, clientId } = params;
  const cookie = request.headers.get("cookie") || "";
  const body = await request.json();

  const [data, error] = await asyncThrowable<ExperienceClient>(() =>
    api.patch<ExperienceClient>(
      clientPath(id ?? "", clientId ?? ""),
      body,
      cookie
    )
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to update client";
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
  const { id, clientId } = params;
  const cookie = request.headers.get("cookie") || "";

  const [_, error] = await asyncThrowable(() =>
    api.delete(clientPath(id ?? "", clientId ?? ""), cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to delete client";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }

  return new Response(JSON.stringify({ data: { message: "Client deleted" } }), {
    status: 200,
    headers: { "Content-Type": "application/json" },
  });
};
