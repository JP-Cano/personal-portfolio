import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { Experience } from "@/types/types";
import { API_PATHS } from "@/utils/constants";

export const GET: APIRoute = async ({ params, request }) => {
  const { id } = params;
  const cookie = request.headers.get("cookie") || "";

  try {
    const data = await api.get<Experience>(
      `${API_PATHS.EXPERIENCES}/${id}`,
      cookie
    );
    return new Response(JSON.stringify({ data }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch experience";
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
    const data = await api.patch<Experience>(
      `${API_PATHS.EXPERIENCES}/${id}`,
      body,
      cookie
    );
    return new Response(JSON.stringify({ data }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to update experience";
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
    await api.delete(`${API_PATHS.EXPERIENCES}/${id}`, cookie);
    return new Response(
      JSON.stringify({ data: { message: "Experience deleted" } }),
      {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }
    );
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to delete experience";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};
