import { useState } from "react";
import UsersSection from "./components/UsersSection";
import UrlsSection from "./components/UrlsSection";
import StatsSection from "./components/StatsSection";

type Tab = "users" | "urls" | "stats";

export default function App() {
  const [tab, setTab] = useState<Tab>("urls");

  return (
    <div className="app">
      <header className="header">
        <h1>🔗 URL Shortener</h1>
        <nav className="nav">
          <button
            className={tab === "urls" ? "active" : ""}
            onClick={() => setTab("urls")}
          >
            Ссылки
          </button>
          <button
            className={tab === "users" ? "active" : ""}
            onClick={() => setTab("users")}
          >
            Пользователи
          </button>
          <button
            className={tab === "stats" ? "active" : ""}
            onClick={() => setTab("stats")}
          >
            Статистика
          </button>
        </nav>
      </header>
      <main className="main">
        {tab === "users" && <UsersSection />}
        {tab === "urls" && <UrlsSection />}
        {tab === "stats" && <StatsSection />}
      </main>
    </div>
  );
}
