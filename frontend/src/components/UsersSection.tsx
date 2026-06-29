import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { User } from "../types";

export default function UsersSection() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // form
  const [nickname, setNickname] = useState("");
  const [email, setEmail] = useState("");

  // edit
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editNickname, setEditNickname] = useState("");
  const [editEmail, setEditEmail] = useState("");

  const load = () => {
    setLoading(true);
    setError("");
    api
      .getUsers()
      .then(setUsers)
      .catch((e) => setError(e.message ?? "Failed to load users"))
      .finally(() => setLoading(false));
  };

  useEffect(load, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    try {
      await api.createUser({ nickname, email });
      setNickname("");
      setEmail("");
      setSuccess("Пользователь создан");
      load();
    } catch (e: any) {
      setError(e.message ?? "Failed to create user");
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Удалить пользователя?")) return;
    setError("");
    try {
      await api.deleteUser(id);
      setSuccess("Пользователь удалён");
      load();
    } catch (e: any) {
      setError(e.message ?? "Failed to delete user");
    }
  };

  const startEdit = (u: User) => {
    setEditingId(u.id);
    setEditNickname(u.nickname);
    setEditEmail(u.email);
  };

  const handlePatch = async (id: number) => {
    setError("");
    setSuccess("");
    try {
      await api.patchUser(id, {
        nickname: editNickname || null,
        email: editEmail || null,
      });
      setEditingId(null);
      setSuccess("Пользователь обновлён");
      load();
    } catch (e: any) {
      setError(e.message ?? "Failed to update user");
    }
  };

  return (
    <div className="card">
      <h2>Пользователи</h2>

      {error && <div className="error" style={{ marginBottom: 12 }}>{error}</div>}
      {success && <div className="success" style={{ marginBottom: 12 }}>{success}</div>}

      <form onSubmit={handleCreate} className="form-row">
        <input
          placeholder="Никнейм"
          value={nickname}
          onChange={(e) => setNickname(e.target.value)}
          required
          minLength={3}
          maxLength={20}
        />
        <input
          placeholder="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <button type="submit">Создать</button>
      </form>

      {loading ? (
        <div className="empty">Загрузка...</div>
      ) : users.length === 0 ? (
        <div className="empty">Нет пользователей</div>
      ) : (
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Никнейм</th>
              <th>Email</th>
              <th>Создан</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id}>
                <td>{u.id}</td>
                <td>
                  {editingId === u.id ? (
                    <input
                      value={editNickname}
                      onChange={(e) => setEditNickname(e.target.value)}
                      style={{ width: 140 }}
                    />
                  ) : (
                    u.nickname
                  )}
                </td>
                <td>
                  {editingId === u.id ? (
                    <input
                      value={editEmail}
                      onChange={(e) => setEditEmail(e.target.value)}
                      style={{ width: 200 }}
                    />
                  ) : (
                    u.email
                  )}
                </td>
                <td>{new Date(u.created_at).toLocaleDateString()}</td>
                <td>
                  <div className="actions">
                    {editingId === u.id ? (
                      <>
                        <button className="btn-sm" onClick={() => handlePatch(u.id)}>
                          Сохранить
                        </button>
                        <button className="btn-sm btn-outline" onClick={() => setEditingId(null)}>
                          Отмена
                        </button>
                      </>
                    ) : (
                      <>
                        <button className="btn-sm btn-outline" onClick={() => startEdit(u)}>
                          ✏️
                        </button>
                        <button className="btn-sm btn-danger" onClick={() => handleDelete(u.id)}>
                          🗑️
                        </button>
                      </>
                    )}
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
