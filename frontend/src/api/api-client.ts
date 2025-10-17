import { PORTFOLIO_BACKEND_URL } from "astro:env/server";
import { ApiError } from "@/types/exceptions.ts";

export class ApiClient {
  private static instance: ApiClient;

  /**
   * Returns the singleton instance of the ApiClient class.
   * Ensures that only one instance of ApiClient exists throughout the application lifecycle.
   *
   * @return {ApiClient} The single instance of the ApiClient.
   */
  static getInstance(): ApiClient {
    if (!ApiClient.instance) {
      ApiClient.instance = new ApiClient();
    }
    return ApiClient.instance;
  }

  /**
   * Sends an HTTP request to a specified endpoint with the given options.
   * Automatically appends base URL, includes credentials, and sets the Content-Type header to JSON.
   *
   * @param {string} endpoint The API endpoint to send the request to.
   * @param {RequestInit} [options] Optional configurations for the HTTP request, such as method, headers, body, etc.
   * @return {Promise<T>} A promise resolving to the parsed JSON response data of type T.
   * @throws {ApiError} Throws an ApiError if the response status is not ok, including the error message or status code.
   */
  async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${PORTFOLIO_BACKEND_URL}/${endpoint}`;

    const response = await fetch(url, {
      ...options,
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new ApiError(error.message || error.error, response.status);
    }

    const result = await response.json();

    return result.data;
  }

  /**
   * Sends a GET request to the specified endpoint and returns the parsed response.
   *
   * @param {string} endpoint - The endpoint URL to send the GET request to.
   * @return {Promise<T>} A promise that resolves to the response of the specified type.
   */
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint);
  }

  /**
   * Sends a POST request to the specified endpoint with the provided data.
   *
   * @param {string} endpoint - The endpoint to which the POST request is sent.
   * @param {unknown} data - The data to include in the body of the POST request.
   * @return {Promise<T>} A promise that resolves to the response of the request as the specified type.
   */
  async post<T>(endpoint: string, data: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  /**
   * Sends an HTTP PATCH request to the specified endpoint with the provided data.
   *
   * @param {string} endpoint The API endpoint to which the PATCH request will be sent.
   * @param {unknown} data The data to be sent in the body of the PATCH request.
   * @return {Promise<T>} A promise that resolves to the server's response data of type T.
   */
  async patch<T>(endpoint: string, data: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PATCH",
      body: JSON.stringify(data),
    });
  }

  /**
   * Sends a DELETE request to the specified endpoint.
   *
   * @param {string} endpoint - The API endpoint to send the DELETE request to.
   * @return {Promise<T>} A promise that resolves with the response data of type T.
   */
  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: "DELETE",
    });
  }
}

export const api = ApiClient.getInstance();
