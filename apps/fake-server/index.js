const express = require("express");
const app = express();
const port = 3022;

const randomInt = (min, max) => {
  return Math.floor(Math.random() * (max - min + 1)) + min;
};


const simpleLogger = (req, res, next) => {
  const start = performance.now();

  res.on('finish', () => {
    const ms = (performance.now() - start).toFixed(1);
    console.log(
      `${req.method} ${req.originalUrl} → ${res.statusCode} (${ms} ms)`
    );
  });

  next();
}

app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(simpleLogger);

const simpleHandler = (req, res, next) => {
  try {
    const delay = parseInt(req.query.delay, 10) || 100;
    const delayRandom = parseInt(req.query.delayRandom, 10) || 20;
    const finalDelay = delay + randomInt(0, delayRandom);

    setTimeout(() => {
      res.send({
        status: "ok",
        delay: finalDelay,
      });
    }, delay);
  } catch (err) {
    next(err);
  }
}

app.get("/", simpleHandler);
app.post("/", simpleHandler);
app.put("/", simpleHandler);
app.delete("/", simpleHandler);
app.patch("/", simpleHandler);

// Global error handler
app.use((err, req, res, next) => {
  console.error('Error:', err);
  res.status(500).json({ error: 'Internal Server Error', details: err.message });
});

// Process-level error handlers
process.on('uncaughtException', (err) => {
  console.error('Uncaught Exception:', err);
});

process.on('unhandledRejection', (reason, promise) => {
  console.error('Unhandled Rejection at:', promise, 'reason:', reason);
});

app.listen(port, () => {
  console.log(`Fake server listening on port ${port}`);
});
