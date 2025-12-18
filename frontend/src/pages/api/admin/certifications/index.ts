import type { APIRoute } from "astro";
import { api } from "@/api/api-client";
import type { CareerCertifications, UploadResponse } from "@/types/types";
import { API_PATHS } from "@/utils/constants";
import { asyncThrowable } from "@/utils/utils.ts";

export const GET: APIRoute = async ({ request }) => {
  const cookie = request.headers.get("cookie") || "";

  const [data, error] = await asyncThrowable<CareerCertifications>(() =>
    api.get<CareerCertifications>(API_PATHS.UPLOAD_CERTIFICATES, cookie)
  );

  if (error) {
    const message =
      error instanceof Error ? error.message : "Failed to fetch certifications";
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
  const formData = await request.formData();

  const [data, error] = await asyncThrowable<UploadResponse>(() =>
    api.postFormData<UploadResponse>(
      API_PATHS.UPLOAD_CERTIFICATES,
      formData,
      cookie
    )
  );

  if (error) {
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

  return new Response(JSON.stringify({ data }), {
    status: 201,
    headers: { "Content-Type": "application/json" },
  });
};
