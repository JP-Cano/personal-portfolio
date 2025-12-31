import { z } from "zod";

export enum WorkType {
  REMOTE = "Remote",
  ON_SITE = "On-site",
  HYBRID = "Hybrid",
}

export const experienceSchema = z.object({
  id: z.number().positive(),
  title: z.string(),
  company: z.string(),
  url: z.url().optional(),
  location: z.string(),
  type: z.enum(Object.values(WorkType)),
  startDate: z.date(),
  endDate: z.date().optional(),
  description: z.string(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

export const projectSchema = z.object({
  id: z.number().positive(),
  name: z.string(),
  description: z.string(),
  url: z.url().optional(),
  startDate: z.date(),
  endDate: z.date().optional(),
  technologies: z.string().optional(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

export const uploadedFileSchema = z.object({
  id: z.number().positive(),
  title: z.string(),
  issuer: z.string(),
  issue_date: z.coerce.date(),
  file_url: z.url(),
  file_name: z.string(),
  original_name: z.string(),
  file_size: z.number(),
  mime_type: z.string(),
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
});

export const uploadErrorSchema = z.object({
  original: z.string(),
  error: z.string(),
});

export const uploadResponseSchema = z.object({
  total: z.number(),
  successful: z.number(),
  failed: z.number(),
  files: z.array(uploadedFileSchema),
  errors: z.array(uploadErrorSchema),
});

export type Experience = z.infer<typeof experienceSchema>;
export type Experiences = Array<Experience>;

export type Project = z.infer<typeof projectSchema>;
export type Projects = Array<Project>;

export type CareerCertifications = Array<z.infer<typeof uploadedFileSchema>>;
export type UploadResponse = z.infer<typeof uploadResponseSchema>;

/**
 * A TypeScript type that represents the result of an asynchronous operation that can either resolve with data
 * or reject with an error. Used to encapsulate both success and failure outcomes in a unified way.
 *
 * The type is defined as a tuple where:
 * - The first element is either the resolved data of type `T` or `null` if an error occurred.
 * - The second element is either the error value (of type `unknown`) or `null` if the operation was successful.
 *
 * This type can help with managing flows where both results and errors need to be captured explicitly.
 *
 * Example structure:
 * - On success: [resolvedData, null]
 * - On failure: [null, errorValue]
 *
 * @template T - The type of data corresponding to the successful resolution of the asynchronous operation.
 */
export type AsyncThrowable<T> =
  | [data: T, error: null]
  | [data: null, error: unknown];
