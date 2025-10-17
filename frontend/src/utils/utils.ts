import type { AsyncThrowable } from "@/types/types.ts";

/**
 * A utility function to handle asynchronous operations with error capturing.
 * It executes a provided asynchronous function and returns an array containing
 * either the resolved result or the captured error.
 *
 * @template T
 * @param {() => Promise<T>} fn - An asynchronous function to execute.
 * @returns {Promise<[T | null, any | null]>} A promise that resolves to
 *          a tuple where the first element is the resolved value or null
 *          if an error occurred, and the second element is the error or null
 *          if it was successful.
 */
export const asyncThrowable = async <T>(
  fn: () => Promise<T>
): Promise<AsyncThrowable<T>> => {
  try {
    const result = await fn();

    return [result, null];
  } catch (err) {
    return [null, err];
  }
};
