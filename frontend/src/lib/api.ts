// Use ?? instead of || so empty string is preserved (for relative URLs via Tailscale)
const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

interface Link {
  title: string;
  url: string;
  type: string;
  description?: string;
}

export interface ProcessedNote {
  id: string;
  original: string;
  title: string;
  category: string;
  markdown: string;
  links: Link[];
  created_at: string;
  synced_at?: string;
}

interface NotesResponse {
  notes: ProcessedNote[];
  total: number;
}

interface CategoryCount {
  name: string;
  count: number;
}

interface CategoriesResponse {
  categories: CategoryCount[];
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: "Unknown error" }));
      throw new Error(error.error || `Request failed: ${response.status}`);
    }

    return response.json();
  }

  async createNote(content: string): Promise<ProcessedNote> {
    return this.request<ProcessedNote>("/api/notes", {
      method: "POST",
      body: JSON.stringify({ content }),
    });
  }

  async getNotes(category?: string, limit = 50, offset = 0): Promise<NotesResponse> {
    const params = new URLSearchParams();
    if (category) params.set("category", category);
    params.set("limit", limit.toString());
    params.set("offset", offset.toString());

    return this.request<NotesResponse>(`/api/notes?${params.toString()}`);
  }

  async getNote(id: string): Promise<ProcessedNote> {
    return this.request<ProcessedNote>(`/api/notes/${id}`);
  }

  async deleteNote(id: string): Promise<void> {
    await this.request(`/api/notes/${id}`, { method: "DELETE" });
  }

  async getCategories(): Promise<CategoriesResponse> {
    return this.request<CategoriesResponse>("/api/categories");
  }

  async healthCheck(): Promise<{ status: string }> {
    return this.request<{ status: string }>("/api/health");
  }
}

export const api = new ApiClient(API_URL);
