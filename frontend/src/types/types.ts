import { z } from "zod";

export enum WorkType {
  REMOTE = "Remote",
  ON_SITE = "On-site",
  HYBRID = "Hybrid",
}

export const experienceSchema = z.object({
  id: z.uint32(),
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
  id: z.uint32(),
  name: z.string(),
  description: z.string(),
  url: z.url().optional(),
  startDate: z.date(),
  endDate: z.date().optional(),
  createdAt: z.date(),
  updatedAt: z.date(),
});

export const experienceRequestSchema = experienceSchema.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export const projectRequestSchema = projectSchema.omit({
  id: true,
  createdAt: true,
  updatedAt: true,
});

export type Experience = z.infer<typeof experienceSchema>;
export type ExperienceRequest = z.infer<typeof experienceRequestSchema>;
export type Experiences = Array<Experience>;

export type Project = z.infer<typeof projectSchema>;
export type ProjectRequest = z.infer<typeof projectRequestSchema>;
export type Projects = Array<Project>;

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
