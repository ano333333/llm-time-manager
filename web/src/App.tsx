import { Route, Routes } from "react-router-dom";
import Capture from "./routes/capture";
import Chat from "./routes/chat";
import Goals from "./routes/goals";
import Home from "./routes/index";
import NotFound from "./routes/not-found";
import Settings from "./routes/settings";
import Tasks from "./routes/tasks";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/chat" element={<Chat />} />
      <Route path="/goals" element={<Goals />} />
      <Route path="/tasks" element={<Tasks />} />
      <Route path="/capture" element={<Capture />} />
      <Route path="/settings/local" element={<Settings />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}

export default App;
