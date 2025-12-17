import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import { API_PATHS } from "@/utils/constants";
import type { Project, Projects } from "@/types/types";

export const GET: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";

  try {
    const data = await api.get<Projects>(API_PATHS.PROJECTS, cookie);
    return new Response(JSON.stringify({ data }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch projects";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};

export const POST: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";
  const body = await request.json();

  try {
    const data = await api.post<Project>(API_PATHS.PROJECTS, body, cookie);
    return new Response(JSON.stringify({ data }), {
      status: 201,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to create project";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};
