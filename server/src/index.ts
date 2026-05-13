import express from 'express';
import type { Request, Response } from 'express';
import cors from 'cors';
import fs from 'fs';
import path from 'path';

const app = express();
const DB_PATH = path.join(__dirname, '../../data/db.json');

app.use(cors());
app.use(express.json());

// Типизируем структуру базы данных
interface DB {
  tournaments: Tournament[];
  teams: Team[];
  matches: Match[];
}

interface Tournament {
  id: string;
  name: string;
  createdAt: string;
}

interface Team {
  id: string;
  name: string;
  tournamentId: string;
}

interface Match {
  id: string;
  team1Id: string;
  team2Id: string;
  winnerId: string | null;
  round: number;
}

const readDB = (): DB => JSON.parse(fs.readFileSync(DB_PATH, 'utf-8'));
const writeDB = (data: DB): void =>
  fs.writeFileSync(DB_PATH, JSON.stringify(data, null, 2));

app.get('/api/ping', (req: Request, res: Response) => {
  res.json({ ok: true, message: 'Сервер работает!' });
});

app.listen(3001, () => {
  console.log('Сервер запущен на http://localhost:3001');
});
