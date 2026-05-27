package main

import (
    "encoding/json"
    "html/template"
    "log"
    "math/rand"
    "net/http"
    "os"
    "strings"
    "time"
)

type botResponse struct {
    Reply string `json:"reply"`
}

type userMessage struct {
    Text string `json:"text"`
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // Словари (те же, что были)
    adviceMap := map[string]string{
        "стол":    "🪵 Отлично! Чтобы добавить старины: сделайте потёртости по краям, покройте матовым лаком с тёмным пигментом.",
        "стул":    "🪑 Совет: состарьте ножки брашированием и нанесите кракелюрный лак на спинку.",
        "лавка":   "🪵 Лавка смотрится с ручной резьбой и царапинами (пройдитесь цепью).",
        "комод":   "🗄️ Для комода подойдёт декупаж с кракелюром или патинирование углов.",
        "полка":   "📚 Полку можно обжечь паяльной лампой и пройтись жёсткой щёткой.",
    }

    styleTips := map[string]string{
        "кантри":    "🌾 Стиль кантри: грубые текстуры, потертости, кованые элементы, пастельные тона.",
        "прованс":   "🌸 Прованс: пастельные оттенки (лавандовый, кремовый), потёртости, цветочные мотивы.",
        "лофт":      "🏭 Лофт: дерево + металл, неровные края, обожжённое дерево, тёмная патина.",
        "шебби-шик": "🎀 Шебби-шик: многослойная покраска, кракелюр, эффект облупившейся краски.",
    }

    randomTips := []string{
        "Перед окрашиванием загрунтуй олифой — трещины станут глубже.",
        "Используй воск с тёмным пигментом для углов и кромок.",
        "Ножки стульев делай точеными — как в начале 20 века.",
        "Готовое изделие припудри мелкой золой для налёта пыли веков.",
        "Не бойся неровностей — ручная работа ценится.",
        "Для эффекта «потертости» пройдись мелкой наждачкой по рёбрам после покраски.",
    }

    botProcess := func(text string) string {
        text = strings.TrimSpace(text)
        lower := strings.ToLower(text)

        switch lower {
        case "/start", "start", "привет", "начать":
            return "🪑 Привет! Я бот для мебели «под старину».\n\nКоманды:\n/ideas — идеи\n/materials — материалы\n/tools — инструменты\n/techniques — техники\n/random_tip — случайный совет\n/gallery — примеры (ссылки)\n/tutorial_stool — урок табурет\n/style — выбрать стиль\n/blueprints — чертежи"
        case "/ideas":
            return "💡 Идеи:\n- «Трапезный» стол из массива сосны\n- Стулья с гнутыми ножками\n- Лавка с резной спинкой\n- Комод с кракелюром\n- Полка с патиной"
        case "/materials":
            return "🌲 Материалы: массив сосны, дуба; фанера + морилка; воск, битумный лак."
        case "/tools":
            return "🔨 Инструменты: рубанок, стамески, шлифмашинка, паяльная лампа, кисти."
        case "/techniques":
            return "🎨 Техники старения: браширование, патинирование, кракелюр, обжиг, царапины."
        case "/random_tip":
            return "✨ " + randomTips[rand.Intn(len(randomTips))]
        case "/gallery":
            return "📸 Примеры (замените ссылки на свои фото):\n- Стол: https://example.com/table1.jpg\n- Стул: https://example.com/chair1.jpg"
        case "/tutorial_stool":
            return "🪑 Пошагово: 1) детали 40x40 см, ножки 45 см; 2) сборка; 3) браширование; 4) обжиг; 5) морилка + воск; 6) потёртости цепью."
        case "/style":
            return "Стили: кантри, прованс, лофт, шебби-шик. Напиши название — расскажу."
        case "/blueprints":
            return "📐 Чертежи: стол 160x80, ножки 74 см; стул 45x45, ножки 45/47 см. Подробнее позже."
        }

        for style, tip := range styleTips {
            if strings.Contains(lower, style) {
                return tip
            }
        }
        for keyword, advice := range adviceMap {
            if strings.Contains(lower, keyword) {
                return advice
            }
        }
        return "Не понял команды 🤔\nНапиши /start для списка команд."
    }

    // HTML-шаблон (тот же красивый чат)
    tmplHTML := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Винтажная мебель — советчик</title>
    <style>
        * { box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Roboto, sans-serif;
            background: #e9e0cf;
            margin: 0;
            padding: 20px;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .chat-container {
            max-width: 800px;
            width: 100%;
            background: #fff7ef;
            border-radius: 24px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            overflow: hidden;
            display: flex;
            flex-direction: column;
            height: 90vh;
        }
        .chat-header {
            background: #5e3a2b;
            color: #f5e7d9;
            padding: 20px;
            font-size: 1.5rem;
            font-weight: bold;
            text-align: center;
        }
        .chat-messages {
            flex: 1;
            overflow-y: auto;
            padding: 20px;
            display: flex;
            flex-direction: column;
            gap: 12px;
        }
        .message {
            max-width: 80%;
            padding: 10px 16px;
            border-radius: 20px;
            line-height: 1.4;
            word-wrap: break-word;
        }
        .user {
            align-self: flex-end;
            background: #d9c5a7;
            color: #2c1a10;
            border-bottom-right-radius: 5px;
        }
        .bot {
            align-self: flex-start;
            background: #f0e3d4;
            color: #2c1a10;
            border-bottom-left-radius: 5px;
        }
        .input-area {
            display: flex;
            padding: 16px;
            background: #faf3e8;
            border-top: 1px solid #e2cfb5;
        }
        #messageInput {
            flex: 1;
            padding: 12px;
            border: 1px solid #d4bca0;
            border-radius: 40px;
            font-size: 1rem;
            outline: none;
        }
        #sendBtn {
            background: #5e3a2b;
            border: none;
            color: white;
            padding: 0 20px;
            margin-left: 10px;
            border-radius: 40px;
            font-size: 1rem;
            font-weight: bold;
            cursor: pointer;
        }
        #sendBtn:hover { background: #7e5240; }
        .footer-note {
            text-align: center;
            font-size: 0.75rem;
            padding: 8px;
            background: #fff7ef;
            color: #7a6a5c;
        }
        @media (max-width: 600px) { .message { max-width: 95%; } }
    </style>
</head>
<body>
<div class="chat-container">
    <div class="chat-header">🪑 Мебельный мастер под старину</div>
    <div class="chat-messages" id="chatMessages">
        <div class="message bot">🪑 Привет! Я помогу советами по мебели «под старину». Напиши /start для списка команд.</div>
    </div>
    <div class="input-area">
        <input type="text" id="messageInput" placeholder="Напишите /ideas, /style или 'стол'...">
        <button id="sendBtn">Отправить</button>
    </div>
    <div class="footer-note">✨ Работает на Go | Доступен из любой точки мира</div>
</div>
<script>
    const messagesDiv = document.getElementById('chatMessages');
    const input = document.getElementById('messageInput');
    const sendBtn = document.getElementById('sendBtn');
    function addMessage(text, isUser) {
        const msgDiv = document.createElement('div');
        msgDiv.className = 'message ' + (isUser ? 'user' : 'bot');
        msgDiv.textContent = text;
        messagesDiv.appendChild(msgDiv);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    }
    async function sendMessage() {
        const text = input.value.trim();
        if (!text) return;
        addMessage(text, true);
        input.value = '';
        try {
            const response = await fetch('/api/message', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({text: text})
            });
            const data = await response.json();
            addMessage(data.reply, false);
        } catch (err) {
            addMessage('❌ Ошибка соединения', false);
        }
    }
    sendBtn.addEventListener('click', sendMessage);
    input.addEventListener('keypress', (e) => { if (e.key === 'Enter') sendMessage(); });
</script>
</body>
</html>`

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        tmpl := template.Must(template.New("chat").Parse(tmplHTML))
        tmpl.Execute(w, nil)
    })

    http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Only POST", http.StatusMethodNotAllowed)
            return
        }
        var msg userMessage
        if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }
        reply := botProcess(msg.Text)
        json.NewEncoder(w).Encode(botResponse{Reply: reply})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Сервер запущен на порту %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}