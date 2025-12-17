import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { Project } from "@/types/types";
import { API_PATHS } from "@/utils/constants";

export const GET: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  try {
    const data = await api.get<Project>(`${API_PATHS.PROJECTS}/${id}`, cookie);
    return new Response(JSON.stringify({ data }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch project";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};

export const PATCH: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";
  const body = await request.json();

  try {
    const data = await api.patch<Project>(
      `${API_PATHS.PROJECTS}/${id}`,
      body,
      cookie
    );
    return new Response(JSON.stringify({ data }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to update project";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};

export const DELETE: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  try {
    await api.delete(`${API_PATHS.PROJECTS}/${id}`, cookie);
    return new Response(
      JSON.stringify({ data: { message: "Project deleted" } }),
      {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }
    );
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to delete project";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};
