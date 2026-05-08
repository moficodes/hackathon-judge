import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fetcher } from './fetcher';

describe('fetcher', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('should fetch and return JSON data on success', async () => {
    const mockData = { id: 1, name: 'Test' };
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockData),
    });

    const result = await fetcher('/api/test');
    expect(result).toEqual(mockData);
    expect(globalThis.fetch).toHaveBeenCalledWith('/api/test');
  });

  it('should throw an error if the response is not ok', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 404,
      statusText: 'Not Found',
    });

    await expect(fetcher('/api/notfound')).rejects.toThrow('An error occurred while fetching the data.');
  });
});
