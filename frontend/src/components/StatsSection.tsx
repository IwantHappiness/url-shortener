import { useState } from "react";
import { api } from "../api/client";
import type { Stats } from "../types";

export default function StatsSection() {
  const [shortUrl, setShortUrl] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!shortUrl.trim()) return;
    setError("");
    setLoading(true);
    try {
      const data = await api.getStats(
        shortUrl.trim(),
        from || undefined,
        to || undefined,
      );
      setStats(data);
    } catch (e: any) {
      setError(e.message ?? "Failed to load stats");
      setStats(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card">
      <h2>Статистика</h2>

      {error && (
        <div className="error" style={{ marginBottom: 12 }}>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="form-row">
        <input
          placeholder="Короткий код (abc123)"
          value={shortUrl}
          onChange={(e) => setShortUrl(e.target.value)}
          required
        />
        <input
          type="date"
          value={from}
          onChange={(e) => setFrom(e.target.value)}
          placeholder="От (ГГГГ-ММ-ДД)"
        />
        <input
          type="date"
          value={to}
          onChange={(e) => setTo(e.target.value)}
          placeholder="До (ГГГГ-ММ-ДД)"
        />
        <button type="submit" disabled={loading}>
          {loading ? "Загрузка..." : "Показать"}
        </button>
      </form>

      {stats && (
        <>
          <table>
            <thead>
              <tr>
                <th>Параметр</th>
                <th>Значение</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Короткий URL</td>
                <td>
                  <a
                    className="link"
                    href={`/${stats.short_url}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {stats.short_url}
                  </a>
                </td>
              </tr>
              <tr>
                <td>Оригинальный URL</td>
                <td>
                  <a
                    className="link"
                    href={stats.original_url}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {stats.original_url}
                  </a>
                </td>
              </tr>
              <tr>
                <td>Создан</td>
                <td>{new Date(stats.created_at).toLocaleString()}</td>
              </tr>
              {stats.last_clicked_at && (
                <tr>
                  <td>Последний переход</td>
                  <td>{new Date(stats.last_clicked_at).toLocaleString()}</td>
                </tr>
              )}
            </tbody>
          </table>

          <div className="stat-grid">
            <div className="stat-card">
              <div className="stat-value">{stats.total_clicks}</div>
              <div className="stat-label">Всего переходов</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">{stats.unique_ips}</div>
              <div className="stat-label">Уникальных IP</div>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
