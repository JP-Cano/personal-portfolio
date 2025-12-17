import { z } from "zod";

export const userSchema = z.object({
  id: z.number().positive(),
  email: z.email(),
});

export const authResponseSchema = z.object({
  user: userSchema,
  message: z.string(),
});

export const loginRequestSchema = z.object({
  email: z.email("Please enter a valid email"),
  password: z.string().min(1, "Password is required"),
});

export type User = z.infer<typeof userSchema>;
export type AuthResponse = z.infer<typeof authResponseSchema>;
export type LoginRequest = z.infer<typeof loginRequestSchema>;
