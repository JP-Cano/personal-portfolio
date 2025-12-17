import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import { API_PATHS } from "@/utils/constants";
import type { CareerCertifications, UploadResponse } from "@/types/types";

export const GET: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";

  try {
    const data = await api.get<CareerCertifications>(
      API_PATHS.UPLOAD_CERTIFICATES,
      cookie
    );
    return new Response(JSON.stringify({ data }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch certifications";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};

export const POST: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";
  const formData = await request.formData();

  try {
    const data = await api.postFormData<UploadResponse>(
      API_PATHS.UPLOAD_CERTIFICATES,
      formData,
      cookie
    );
    return new Response(JSON.stringify({ data }), {
      status: 201,
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Failed to upload certifications";
    const status = (error as { statusCode?: number }).statusCode || 500;
    return new Response(JSON.stringify({ error: message }), {
      status,
      headers: { "Content-Type": "application/json" },
    });
  }
};
