/// <reference types="astro/client" />

import type { User } from "@/types/auth";

declare global {
  namespace App {
    interface Locals {
      user?: User;
    }
  }
}

export {};
