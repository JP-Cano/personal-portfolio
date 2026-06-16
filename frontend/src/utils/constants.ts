export const API_PATHS = {
  EXPERIENCES: "experiences",
  // Nested sub-resource: clients that belong to a given experience.
  EXPERIENCE_CLIENTS: (experienceId: number | string) =>
    `experiences/${experienceId}/clients`,
  PROJECTS: "projects",
  UPLOAD_CERTIFICATES: "upload-certificates",
  AUTH: {
    LOGIN: "auth/login",
    LOGOUT: "auth/logout",
    ME: "auth/me",
  },
};
