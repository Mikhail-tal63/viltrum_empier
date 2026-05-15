الفكرة الأساسية:

بدل ما الفرونت كل مرة يعمل:

```txt
GET /tasks
```

ويعمل refresh…

حنفتح اتصال دائم بين الفرونت والباک:

# WebSocket

الاتصال يظل مفتوح طول ما المستخدم داخل التطبيق.

---

# كيف الريلتايم حيشتغل فعليًا

عندك 3 أجزاء:

```txt
Frontend
↕ websocket
Go Server
↕
Redis Pub/Sub
```

---

# السيناريو الكامل

## أحمد يحرك Task

من:

```txt
Todo
```

إلى:

```txt
In Progress
```

---

# الخطوة 1 — Frontend يرسل Event

الفرونت يرسل عبر websocket:

```json
{
  "type": "task_moved",
  "task_id": "123",
  "from_column": "todo",
  "to_column": "progress"
}
```

---

# الخطوة 2 — Go يستقبل الحدث

الـ websocket handler يستقبل الرسالة.

---

# الخطوة 3 — تحدث MongoDB

السيرفر يعدل:

```txt
column_id
position
updated_at
```

في task.

---

# الخطوة 4 — Publish في Redis

بعد نجاح التعديل:

```txt
workspace:1:events
```

ينرسل له event:

```json
{
  "type": "task_updated",
  "task": {...}
}
```

---

# الخطوة 5 — كل السيرفرات تستقبل

أي websocket server مشترك في Redis channel
يستقبل الرسالة.

---

# الخطوة 6 — Broadcast للكل

الـ Hub يرسل الحدث لكل المستخدمين داخل نفس الـ board/workspace.

---

# النتيجة

كل الناس تشوف الكارد تتحرك فورًا.

بدون refresh.

---

# ليش Redis مهم؟

لو عندك websocket server واحد فقط:
تقدر تبعث مباشرة.

لكن لو لاحقًا صار عندك:

* server 1
* server 2
* server 3

كيف server 2 يعرف إن server 1 حدث التاسك؟

هنا Redis Pub/Sub يحل المشكلة.

---

# WebSocket Architecture

كل مستخدم عنده:

```go
type Client struct {
    Conn *websocket.Conn
    Send chan []byte
}
```

---

# والـ Hub يدير الاتصالات

```go
type Hub struct {
    Clients map[*Client]bool
    Broadcast chan []byte
    Register chan *Client
    Unregister chan *Client
}
```

---

# الريلتايم اللي عندك

## 1. Task Updates

* create
* edit
* move
* delete

---

## 2. Comments

أحد كتب comment:
الكل يشوفه فورًا.

---

## 3. Notifications

```txt
Ahmed assigned you to task
```

توصل مباشرة.

---

## 4. Online Presence

لما المستخدم يتصل:

* يتحط online في Redis
* يتبعت event

الكل يشوف:
🟢 online

---

# كيف تعرف الشخص online؟

أول ما websocket يتصل:

```txt
online:user:15 = true
```

في Redis.

مع expiration.

ولو disconnect:
يتشال.

---

# Typing Indicator

لو المستخدم يكتب comment:

frontend يرسل:

```json
{
  "type":"typing",
  "task_id":"123"
}
```

السيرفر يبعتها للباقي.

فيطلع:

```txt
Ahmed is typing...
```

---

# Activity Feed

أي حدث:

* task moved
* comment added
* member joined

ينرسل realtime.

---

# هل كل شيء لازم websocket؟

لا.

---

# REST API

تستخدمه لـ:

* login
* initial fetch
* CRUD الطبيعي

---

# WebSocket

تستخدمه لـ:

* live updates
* notifications
* online users
* typing
* realtime sync

---

# الشكل الحقيقي للتطبيق

أول ما تدخل board:

## الفرونت:

يجيب البيانات الأساسية عبر REST.

---

بعدها:
يفتح websocket.

---

من هنا أي تغيير:
يجي realtime فقط.

---

# أهم نقطة لازم تفهمها

WebSocket ليس Database.

هو فقط:

# قناة أحداث live

أما الحقيقة الأساسية للبيانات:
تبقى MongoDB.

---

# ليش هذا التصميم ممتاز؟

لأن:

* MongoDB = source of truth
* Redis = realtime broker/cache
* WebSocket = live transport

وهذا تقريبًا نفس الفكر الموجود في تطبيقات حقيقية كثيرة.
