import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { User, Url } from "../types";

export default function UrlsSection() {
  const [urls, setUrls] = useState<Url[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // form
  const [originalUrl, setOriginalUrl] = useState("");
  const [userId, setUserId] = useState("");

  // filter
  const [filterUserId, setFilterUserId] = useState("");

  const loadUsers = () =>
    api
      .getUsers()
      .then(setUsers)
      .catch(() => {});

  const load = (uid?: number) => {
    setLoading(true);
    setError("");
    api
      .getURLs(uid)
      .then(setUrls)
      .catch((e) => setError(e.message ?? "Failed to load URLs"))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    loadUsers();
    load();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    try {
      await api.createURL({ url: originalUrl, user_id: Number(userId) });
      setOriginalUrl("");
      setSuccess("Ссылка создана");
      load(filterUserId ? Number(filterUserId) : undefined);
    } catch (e: any) {
      setError(e.message ?? "Failed to create URL");
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Удалить ссылку?")) return;
    setError("");
    try {
      await api.deleteURL(id);
      setSuccess("Ссылка удалена");
      load(filterUserId ? Number(filterUserId) : undefined);
    } catch (e: any) {
      setError(e.message ?? "Failed to delete URL");
    }
  };

  const handleRegen = async (id: number) => {
    setError("");
    setSuccess("");
    try {
      await api.patchURL(id);
      setSuccess("Короткая ссылка обновлена");
      load(filterUserId ? Number(filterUserId) : undefined);
    } catch (e: any) {
      setError(e.message ?? "Failed to regenerate URL");
    }
  };

  const handleFilter = () => {
    load(filterUserId ? Number(filterUserId) : undefined);
  };

  const copyToClipboard = (shortUrl: string) => {
    const fullUrl = `${window.location.origin}/${shortUrl}`;
    navigator.clipboard.writeText(fullUrl).then(() => {
      setSuccess(`Скопировано: ${fullUrl}`);
      setTimeout(() => setSuccess(""), 2000);
    });
  };

  return (
    <div className="card">
      <h2>Ссылки</h2>

      {error && (
        <div className="error" style={{ marginBottom: 12 }}>
          {error}
        </div>
      )}
      {success && (
        <div className="success" style={{ marginBottom: 12 }}>
          {success}
        </div>
      )}

      <form onSubmit={handleCreate} className="form-row">
        <input
          placeholder="https://example.com"
          value={originalUrl}
          onChange={(e) => setOriginalUrl(e.target.value)}
          required
        />
        <select
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          required
        >
          <option value="">Выберите пользователя</option>
          {users.map((u) => (
            <option key={u.id} value={u.id}>
              {u.nickname} ({u.id})
            </option>
          ))}
        </select>
        <button type="submit">Создать</button>
      </form>

      <div className="form-row" style={{ marginTop: 8 }}>
        <select
          value={filterUserId}
          onChange={(e) => {
            setFilterUserId(e.target.value);
          }}
        >
          <option value="">Все пользователи</option>
          {users.map((u) => (
            <option key={u.id} value={u.id}>
              {u.nickname} ({u.id})
            </option>
          ))}
        </select>
        <button className="btn-outline" onClick={handleFilter}>
          Фильтр
        </button>
      </div>

      {loading ? (
        <div className="empty">Загрузка...</div>
      ) : urls.length === 0 ? (
        <div className="empty">Нет ссылок</div>
      ) : (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Оригинал</th>
              <th>Короткая</th>
              <th>Владелец</th>
              <th>Создан</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {urls.map((u) => (
              <tr key={u.id}>
                <td>{u.id}</td>
                <td
                  style={{
                    maxWidth: 200,
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                >
                  <a
                    className="link"
                    href={u.original_url}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {u.original_url}
                  </a>
                </td>
                <td>
                  <a
                    className="short-url-copy link"
                    href={`/${u.short_url}`}
                    target="_blank"
                    rel="noreferrer"
                    onClick={(e) => {
                      e.stopPropagation();
                      copyToClipboard(u.short_url);
                    }}
                    title="Открыть редирект"
                  >
                    {u.short_url} 📋
                  </a>
                </td>
                <td>
                  {users.find((x) => x.id === u.user_id)?.nickname ?? u.user_id}
                </td>
                <td>{new Date(u.created_at).toLocaleDateString()}</td>
                <td>
                  <div className="actions">
                    <button
                      className="btn-sm btn-outline"
                      onClick={() => handleRegen(u.id)}
                      title="Сгенерировать новый код"
                    >
                      🔄
                    </button>
                    <button
                      className="btn-sm btn-danger"
                      onClick={() => handleDelete(u.id)}
                    >
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
