import { PORTFOLIO_BACKEND_URL } from "astro:env/server";
import type { AuthResponse, LoginRequest, User } from "@/types/auth";
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
   * Login with email and password.
   * The backend will set an HttpOnly cookie with the session ID.
   * Returns both the auth response and the Set-Cookie header for SSR forwarding.
   */
  async login(
    credentials: LoginRequest
  ): Promise<{ data: AuthResponse; setCookie: string | null }> {
    const url = `${PORTFOLIO_BACKEND_URL}/auth/login`;

    const response = await fetch(url, {
      method: "POST",
      body: JSON.stringify(credentials),
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new ApiError(error.message || error.error, response.status);
    }

    const result = await response.json();
    const setCookie = response.headers.get("set-cookie");

    return { data: result.data, setCookie };
  }

  /**
   * Get current authenticated user.
   * Returns user data if session is valid, throws if not authenticated.
   * @param cookie - Optional cookie header for SSR requests
   */
  async getCurrentUser(cookie?: string): Promise<User> {
    const headers: HeadersInit = {};
    if (cookie) {
      headers["Cookie"] = cookie;
    }
    return this.request<User>("auth/me", { headers });
  }

  /**
   * Logout the current user.
   * The backend will clear the session cookie.
   * @param cookie - Optional cookie header for SSR requests
   */
  async logout(cookie?: string): Promise<void> {
    const headers: HeadersInit = {};
    if (cookie) {
      headers["Cookie"] = cookie;
    }
    await this.request<{ message: string }>("auth/logout", {
      method: "POST",
      headers,
    });
  }

  /**
   * Builds headers object with optional cookie for SSR requests.
   * @param cookie - Optional cookie header for SSR requests
   */
  private buildHeaders(cookie?: string): HeadersInit {
    const headers: HeadersInit = {};
    if (cookie) {
      headers["Cookie"] = cookie;
    }
    return headers;
  }

  /**
   * Sends a GET request to the specified endpoint and returns the parsed response.
   *
   * @param {string} endpoint - The endpoint URL to send the GET request to.
   * @param {string} [cookie] - Optional cookie header for SSR requests.
   * @return {Promise<T>} A promise that resolves to the response of the specified type.
   */
  async get<T>(endpoint: string, cookie?: string): Promise<T> {
    return this.request<T>(endpoint, { headers: this.buildHeaders(cookie) });
  }

  /**
   * Sends a POST request to the specified endpoint with the provided data.
   *
   * @param {string} endpoint - The endpoint to which the POST request is sent.
   * @param {unknown} data - The data to include in the body of the POST request.
   * @param {string} [cookie] - Optional cookie header for SSR requests.
   * @return {Promise<T>} A promise that resolves to the response of the request as the specified type.
   */
  async post<T>(endpoint: string, data: unknown, cookie?: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: JSON.stringify(data),
      headers: this.buildHeaders(cookie),
    });
  }

  /**
   * Sends an HTTP PATCH request to the specified endpoint with the provided data.
   *
   * @param {string} endpoint The API endpoint to which the PATCH request will be sent.
   * @param {unknown} data The data to be sent in the body of the PATCH request.
   * @param {string} [cookie] - Optional cookie header for SSR requests.
   * @return {Promise<T>} A promise that resolves to the server's response data of type T.
   */
  async patch<T>(endpoint: string, data: unknown, cookie?: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PATCH",
      body: JSON.stringify(data),
      headers: this.buildHeaders(cookie),
    });
  }

  /**
   * Sends a DELETE request to the specified endpoint.
   *
   * @param {string} endpoint - The API endpoint to send the DELETE request to.
   * @param {string} [cookie] - Optional cookie header for SSR requests.
   * @return {Promise<T>} A promise that resolves with the response data of type T.
   */
  async delete<T>(endpoint: string, cookie?: string): Promise<T> {
    return this.request<T>(endpoint, {
      method: "DELETE",
      headers: this.buildHeaders(cookie),
    });
  }

  /**
   * Sends a POST request with FormData (for file uploads).
   *
   * @param {string} endpoint - The endpoint to which the POST request is sent.
   * @param {FormData} formData - The form data to send.
   * @param {string} [cookie] - Optional cookie header for SSR requests.
   * @return {Promise<T>} A promise that resolves to the response of the request.
   */
  async postFormData<T>(
    endpoint: string,
    formData: FormData,
    cookie?: string
  ): Promise<T> {
    const url = `${PORTFOLIO_BACKEND_URL}/${endpoint}`;
    const headers: HeadersInit = this.buildHeaders(cookie);

    const response = await fetch(url, {
      method: "POST",
      body: formData,
      headers,
      credentials: "include",
    });

    if (!response.ok) {
      const error = await response.json();
      throw new ApiError(error.message || error.error, response.status);
    }

    const result = await response.json();
    return result.data;
  }
}

export const api = ApiClient.getInstance();
