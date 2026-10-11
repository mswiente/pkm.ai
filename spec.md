# AI-basiertes Personal Knowledge System

> **Stand 2026-10-11:** Die Abschnitte 1–17 beschreiben den bisherigen
> CLI-/Obsidian-Ansatz einschließlich seiner damaligen Ausbauideen. Der neue,
> noch nicht implementierte Zielstand steht in [Abschnitt 18: Geplant](#18-geplant-github-markdown-und-mcp).
> Für die Neuentwicklung hat Abschnitt 18 bei Abweichungen Vorrang; die bestehende
> CLI bleibt bis zur ausdrücklich durchgeführten Migration unverändert.

## 1. Zielbild

Wir bauen ein persönliches Knowledge System auf Basis von **Obsidian** als Wissensoberfläche und **Claude Code** als AI-gestütztem Arbeitswerkzeug.

Das System soll:

* Wissen langfristig strukturiert speichern
* Notizen, Projekte, Aufgaben und Referenzen miteinander verknüpfen
* AI bei der Erstellung, Verdichtung, Suche, Umstrukturierung und Weiterverarbeitung unterstützen
* möglichst lokal, nachvollziehbar und markdown-basiert bleiben
* schrittweise ausbaubar sein, ohne früh zu viel Komplexität einzuführen

## 2. Scope

### Im Scope

* Obsidian Vault als zentrale Wissensbasis
* Markdown-first Dokumentation
* Claude Code als Assistenz für:

  * Erfassen und Umformulieren von Notizen
  * Zusammenfassen von Inhalten
  * Strukturieren und Verlinken
  * Ableiten von To-dos, offenen Fragen und nächsten Schritten
  * Erstellen und Pflegen von Templates und Regeln
* klare Ordnerstruktur
* Konventionen für Dateinamen, Metadaten und Links
* einfacher Workflow für Inbox -> Verarbeitung -> Wissensbasis
* optionale spätere Automatisierung

### Nicht im ersten Schritt

* komplexe Multi-Agenten-Architektur
* vollautomatische autonome Wissenspflege
* aufwendige Vektor-Datenbank oder eigener RAG-Stack
* mobile-first Speziallösungen
* tiefe Integration in viele Fremdsysteme

## 3. Designprinzipien

1. **Markdown als Source of Truth**
   Alle inhaltlichen Assets liegen als lesbare Markdown-Dateien im Vault.

2. **Human-in-the-loop**
   AI unterstützt, entscheidet aber nicht autonom über dauerhafte Wissensstruktur.

3. **Einfach vor vollständig**
   Erst eine robuste Basis schaffen, dann Automatisierung ergänzen.

4. **Explizite Struktur plus flexible Verlinkung**
   Wenige klare Regeln, aber keine übermäßige Bürokratie.

5. **Nachvollziehbarkeit vor Magie**
   Änderungen durch AI sollen überprüfbar und verständlich sein.

6. **Portable und tool-unabhängige Inhalte**
   Das System soll auch ohne AI und notfalls ohne Obsidian weiter nutzbar bleiben.

## 4. Haupt-Use-Cases

### 4.1 Export aus ChatGPT nach Obsidian

Wenn Inhalte in ChatGPT entstehen, soll der Nutzer bei Bedarf eine strukturierte Zusammenfassung als Markdown erzeugen können, um diese in Obsidian zu übernehmen.

Wichtige Regel:
Alle erzeugten Markdown-Dateien werden zunächst in `00-inbox/` abgelegt.

Erst in einem separaten AI-gestützten Processing-Schritt werden sie:

* klassifiziert
* erweitert
* verlinkt
* und in Zielordner verschoben

Ziele:

* relevante Antworten verdichten statt komplette Chats ungefiltert abzulegen
* saubere Markdown-Struktur für spätere Weiterverarbeitung
* spätere automatische Klassifikation ermöglichen
* manuelle Qualitätskontrolle vor finaler Ablage

Typische Inhalte:

* Zusammenfassungen technischer Diskussionen
* Architekturüberlegungen
* Entscheidungsoptionen mit Pros/Cons
* Schritt-für-Schritt-Anleitungen
* Entwürfe für Dokumente oder Mails

### 4.2 Export aus Coding-Agent-Sessions nach Obsidian

Wenn mit einem Coding Agent gearbeitet wird, soll die Session bei Bedarf in eine persistente Wissensnotiz überführt werden.

Wichtige Regel:
Auch diese Markdown-Dateien landen zunächst immer in `00-inbox/`.

Erst danach erfolgt ein strukturierender AI-Processing-Schritt.

Ziele:

* Fehlerbehebungen nachvollziehbar dokumentieren
* Architekturentscheidungen aus Implementierungsarbeit extrahieren
* wichtige Learnings aus Debugging-Sessions sichern
* keine vollständigen Roh-Logs ablegen, sondern kuratierte Markdown-Artefakte

Typische Inhalte:

* Bugfix-Notizen
* Root-Cause-Analysen
* Troubleshooting Playbooks
* Architektur- oder Designentscheidungen
* Implementierungsnotizen zu APIs, Deployments oder Konfigurationen

### 4.3 AI-gestützte Erstellung von Standardnotizen

Claude Code soll beim Erstellen und Pflegen standardisierter Notizen unterstützen.

Im Fokus stehen:

* Daily Notes
* Meeting Notes
* Decision Notes
* Project Notes
* Knowledge Notes

Ziele:

* einheitliche Templates nutzen
* schnelle Erstbefüllung aus Stichpunkten oder Rohtext
* konsistente Metadaten und Links erzeugen
* Nachbearbeitung minimieren

### 4.4 Unterstützung mehrerer AI-CLI-Werkzeuge

Das System soll nicht nur mit Claude Code funktionieren, sondern auch mit GitHub Copilot CLI.

Ziele:

* tool-unabhängige Markdown-Artefakte
* gemeinsame Konventionen für Exporte, Templates und Ordnerstruktur
* austauschbare Workflows, soweit möglich
* keine zu starke Kopplung an einen einzelnen Anbieter

### 4.5 Capture

Schnelles Erfassen von:

* Ideen
* Meeting-Notizen
* Projektgedanken
* Learnings
* Links / Quellen
* offene Fragen

### 4.6 Organize

* Notizen bereinigen
* umbenennen
* verschlagworten
* in passende Bereiche verschieben
* mit bestehenden Notizen verlinken

### 4.7 Distill

* lange Notizen verdichten
* Kernaussagen extrahieren
* Entscheidungen, Risiken und To-dos ableiten
* Executive Summary erzeugen

### 4.8 Retrieve

* relevante frühere Notizen finden
* thematische Zusammenhänge erkennen
* verwandte Projekte, Konzepte und Entscheidungen zusammenführen

### 4.9 Create

* aus vorhandenem Wissen neue Artefakte erstellen:

  * Memos
  * Architektur-Skizzen
  * Entscheidungsdokumente
  * E-Mails
  * Präsentationsentwürfe

## 5. Zielarchitektur

## 5.0 Capture Interfaces

Das System unterstützt mehrere standardisierte Wege zur Erfassung neuer Inhalte in den Vault. Alle Capture-Wege erzeugen zunächst Markdown-Dateien in `00-inbox/`.

Unterstützte Capture-Interfaces:

1. Obsidian (manuelle Eingabe)
2. Capture CLI (`pkm-capture`)
3. ChatGPT Markdown-Export
4. Claude Code Session-Export
5. GitHub Copilot CLI Session-Export

Alle Interfaces folgen denselben Regeln:

* erzeugen Markdown
* schreiben nach `00-inbox/`
* enthalten minimales Frontmatter
* sind AI-processing-ready

Das Ziel ist ein einheitlicher Intake-Layer für das gesamte Personal Knowledge System.

## 5.1 Kernkomponenten

### Obsidian

Verwaltet den Vault als primäre Benutzeroberfläche für:

* Lesen
* Schreiben
* Verlinken
* Graph-/Backlink-Navigation
* manuelle Wissenspflege

### Claude Code

Wird auf dem lokalen Vault oder ausgewählten Teilbereichen eingesetzt für:

* Bearbeitung von Markdown-Dateien
* Erzeugen neuer Notizen aus Vorlagen
* Refactoring von Wissensstrukturen
* Zusammenfassungen und Umstrukturierungen
* halbautomatische Pflege durch definierte Prompts / Workflows
* Erstellen standardisierter Notizen wie Daily, Meeting und Decision Notes

### GitHub Copilot CLI

Wird als alternatives oder ergänzendes CLI-Werkzeug unterstützt für:

* Zusammenfassen von Sessions
* Erzeugen von Markdown-Notizen aus Entwicklungsarbeit
* Ableiten von Bugfix-Notizen, Learnings und Architekturentscheidungen
* Bearbeiten von Notizen im Vault über definierte Workflows

### Git

Optional, aber empfohlen für:

* Versionierung
* Nachvollziehbarkeit von Änderungen
* Experimentieren mit AI-gestützten Refactorings
* Backup / Synchronisation

## 5.2 Arbeitsmodell

Der Nutzer arbeitet primär in Obsidian.
Claude Code arbeitet auf dem Dateisystem bzw. Repository und unterstützt gezielt bei Wissensarbeit.

Typischer Zyklus:

1. Inhalt erfassen
2. In Inbox oder Tagesnotiz ablegen
3. Mit Claude Code aufbereiten
4. In dauerhafte Struktur überführen
5. Später wiederfinden und weiterverwenden

## 6. Capture CLI

Das CLI-Tool heißt **`pkm`**. Das zugehörige Repository heißt **`pkm.ai`**.

Ziel:
Ein leichtgewichtiges, erweiterbares Command-Line-Interface für Capture, Verarbeitung und Pflege des Personal Knowledge Systems.

Designprinzip:

* `pkm` ist das einheitliche Frontend für lokale PKM-Workflows
* Markdown-Dateien bleiben die Source of Truth
* Capture und Curation sind bewusst getrennt
* der MVP fokussiert auf schnelle Erfassung in `00-inbox/`

### 6.1 CLI-Kommandostruktur

Grundform:

```bash
pkm <command> [subcommand] [options]
```

Implementierte Kommandos (Stand 2026-04):

```bash
pkm capture [text] [--title] [--source] [--tags] [--type-hint] [--editor] [--clipboard] [--update]
pkm process inbox [--file] [--all] [--full] [--dry-run] [--apply] [--interactive]
pkm daily create [--date] [--open]
pkm meeting create [--title] [--date] [--project] [--participants]
pkm decision create [--title] [--project] [--status] [--from-stdin]
pkm note move <filename> <folder> [--type] [--status] [--dry-run]
pkm knowledge append-topic <slug> --title <title> [--dry-run]
pkm knowledge update-index <slug> --description <desc> [--dry-run]
pkm knowledge append-log --note <file> --action <action> [--filed-to] [--updated] [--created]
pkm project update <slug> [--title] [--intent] [--current-status] [--next-steps] [--plan-heading] [--status] [--dry-run]
pkm project list
pkm sync readwise [--dry-run] [--since] [--limit]
pkm sync readwise auth
pkm config --show | --set-vault-path <path>
pkm skill list | install [name]
```

## 6.2 MVP-Kommandos

Im MVP sind folgende Kommandos vorgesehen:

### `pkm capture`

Erfasst eine neue Inbox-Note.

Beispiele:

```bash
pkm capture "Idee für Workshop zu AI Knowledge Systems"
```

```bash
pbpaste | pkm capture
```

```bash
cat exported-chat.md | pkm capture --source chatgpt --title "RAG vs PKM"
```

Verhalten:

* schreibt immer nach `00-inbox/`
* erzeugt eine Markdown-Datei
* setzt minimales Frontmatter
* akzeptiert Textargument oder stdin
* vergibt Dateinamen nach definierter Konvention
* führt keine Klassifikation oder Verschiebung durch

Optionale Flags:

* `--title <text>`
* `--source <manual|chatgpt|claude-code|copilot-cli|other>`
* `--tags <csv>`
* `--type-hint <value>`
* `--editor` öffnet vor dem Speichern einen Editor
* `--clipboard` liest Inhalt direkt aus der Zwischenablage

### `pkm process inbox`

Startet den AI-gestützten Inbox-Processing-Workflow.

Ziel:
Inbox-Notes analysieren, anreichern und für die Verschiebung vorbereiten.

Beispiele:

```bash
pkm process inbox
```

```bash
pkm process inbox --file 2026-04-03-1842-chatgpt-rag-vs-pkm.md
```

```bash
pkm process inbox --dry-run
```

Verhalten:

* analysiert eine oder mehrere Dateien in `00-inbox/`
* erkennt wahrscheinlichen Notiztyp
* verbessert Titel
* ergänzt Frontmatter
* schlägt Tags und Links vor
* schlägt Zielordner vor
* erzeugt optional eine Vorschau der Änderungen
* verschiebt Dateien nur nach Bestätigung oder expliziter Freigabe

Optionale Flags:

* `--file <filename>`
* `--all`
* `--dry-run`
* `--apply`
* `--interactive`

### `pkm daily create`

Erzeugt eine neue Daily Note aus dem Template.

Beispiele:

```bash
pkm daily create
```

```bash
pkm daily create --date 2026-04-03
```

Verhalten:

* erzeugt eine Datei unter `01-daily/`
* nutzt das Daily-Template
* legt Datei nur an, wenn sie noch nicht existiert

Optionale Flags:

* `--date <YYYY-MM-DD>`
* `--open`

### `pkm meeting create`

Erzeugt eine Meeting Note aus Template oder Rohinput.

Beispiele:

```bash
pkm meeting create --title "Jour fixe Platform Team"
```

```bash
pbpaste | pkm meeting create --title "Architecture Sync"
```

Verhalten:

* erzeugt eine Meeting Note zunächst in `00-inbox/` oder optional direkt in einem Projektkontext
* nutzt ein standardisiertes Template
* kann Rohinput übernehmen und vorstrukturieren

Optionale Flags:

* `--title <text>`
* `--date <YYYY-MM-DD>`
* `--project <name>`
* `--participants <csv>`
* `--inbox`

### `pkm decision create`

Erzeugt eine Decision Note aus Template.

Beispiele:

```bash
pkm decision create --title "Use Claude Code and Copilot CLI as supported AI tools"
```

Verhalten:

* erzeugt eine Decision Note aus Template
* kann aus leerem Template oder aus Rohtext starten
* wird standardmäßig in `00-inbox/` erzeugt, wenn die Entscheidung noch kuratiert werden soll

Optionale Flags:

* `--title <text>`
* `--project <name>`
* `--status <draft|accepted|superseded>`
* `--from-stdin`
* `--inbox`

## 6.3 Spätere Kommandos

Noch nicht implementiert, aber vorgesehen:

### `pkm note refine`

Verbessert eine bestehende Notiz AI-gestützt.

### `pkm note link`

Schlägt interne Links und verwandte Notizen vor.

### `pkm export chat`

Standardisiert Export aus ChatGPT oder anderen Chat-Tools.

### `pkm config init`

Initialisiert `pkm` in einem bestehenden oder neuen Vault.

## 6.4 Dateinamen-Konvention

Für Inbox-Notes wird folgendes Schema empfohlen:

```text
YYYY-MM-DD-HHMM-source-slug.md
```

Beispiele:

```text
2026-04-03-1842-chatgpt-rag-vs-pkm.md
2026-04-03-1905-claude-code-alb-oidc-debugging.md
2026-04-03-1912-manual-idea-workshop-ai-knowledge-system.md
```

Ziele:

* chronologische Sortierung
* Herkunft sofort sichtbar
* stabile Verarbeitung im CLI und in Scripts

Für Daily Notes:

```text
YYYY-MM-DD.md
```

Für andere strukturierte Notizen kann später eine eigene Konvention ergänzt werden.

## 6.5 Minimales Dateiformat für `pkm capture`

Beispiel:

```yaml
---
title: Fix ALB OIDC Redirect Issue
type: inbox
status: inbox
source: copilot-cli
tags: [aws, oidc, alb]
created: 2026-04-03
updated: 2026-04-03
---

## Context

Kurzbeschreibung oder automatisch übernommener Inhalt

## Content

<captured text>
```

Pflichtfelder im MVP:

* `title`
* `type`
* `status`
* `source`
* `created`

Optionale Felder:

* `tags`
* `updated`
* `type_hint`

## 6.6 Interaktionsmodell

`pkm` soll drei Modi unterstützen:

### Non-interactive

Für Scripts und Pipes.

Beispiel:

```bash
pbpaste | pkm capture --source chatgpt --title "Prompt Summary"
```

### Interactive

Für terminalgeführte Auswahl und Bestätigung.

Beispiel:

```bash
pkm process inbox --interactive
```

### Editor-based

Für manuelle Ergänzung im bevorzugten Editor.

Beispiel:

```bash
pkm capture --editor
```

## 6.7 Konfiguration

`pkm` benötigt eine lokale Konfiguration, z. B. über eine Datei wie:

```text
~/.config/pkm/config.yaml
```

Beispielinhalte:

```yaml
vault_path: ~/Documents/Obsidian/pkm-vault
inbox_path: 00-inbox
daily_path: 01-daily
templates_path: 07-templates
editor: code
filename_timezone: Europe/Berlin
default_source: manual
```

## 6.8 Fehler- und Sicherheitsprinzipien

* keine Datei wird stillschweigend außerhalb des Vaults erzeugt
* keine Inbox-Datei wird ohne explizite Freigabe verschoben, außer der Nutzer konfiguriert dies bewusst
* `--dry-run` soll für verändernde AI-Prozesse verfügbar sein
* AI-Vorschläge sollen nachvollziehbar bleiben
* Rohinhalte werden vor einer Umstrukturierung nicht gelöscht

## 6.9 Integration mit AI-Tools

`pkm` ist optimiert für Nutzung mit:

* Claude Code
* GitHub Copilot CLI
* ChatGPT Exports
* Shell Scripts
* Raycast / Alfred / Shortcuts

Damit entsteht ein standardisierter Intake- und Processing-Mechanismus für alle AI-generierten Inhalte.

## 7. Informationsarchitektur

## 6.1 Vorgeschlagene Ordnerstruktur

```text
vault/
  00-inbox/
  01-daily/
  02-projects/
  03-areas/
  04-knowledge/
  05-resources/
  06-decisions/
  07-templates/
  08-attachments/
  09-archive/
```

## 6.2 Bedeutung der Bereiche

### 00-inbox

Unsortierte schnelle Eingaben, die später verarbeitet werden.

### 01-daily

Tagesnotizen als leichtgewichtiger Einstiegspunkt für Gedanken, Meetings und Logbuch.

### 02-projects

Zeitlich oder zielbezogene Vorhaben mit klarem Outcome.

### 03-areas

Dauerhafte Verantwortungsbereiche, z. B. Karriere, Gesundheit, Finanzen, Lernen.

### 04-knowledge

Evergreen Notes, Konzepte, Modelle, Zusammenhänge.

### 05-resources

Quellen, Referenzen, Ausschnitte, Literatur- und Linknotizen.

### 06-decisions

Architektur- oder Arbeitsentscheidungen, inklusive Kontext und Begründung.

### 07-templates

Vorlagen für wiederkehrende Notiztypen.

### 08-attachments

Bilder, PDFs, Exporte, andere Dateien.

### 09-archive

Abgeschlossene oder veraltete Inhalte.

## 7. Notiztypen

## 7.1 Inbox Note

Für rohe Eingaben ohne Anspruch auf Struktur.

## 7.2 Daily Note

Zeitlicher Container für Aktivitäten, Gedanken, Erkenntnisse und Rohnotizen.

## 7.3 Project Note

Enthält Ziel, Kontext, Status, nächste Schritte, relevante Links und Artefakte.

## 7.4 Knowledge Note

Verdichtete, wiederverwendbare Erkenntnis zu einem Thema.

## 7.5 Source / Resource Note

Quelle mit Kurzinhalt, Relevanz und Verweisen auf abgeleitete Erkenntnisse.

## 7.6 Decision Note

Dokumentiert Entscheidung, Optionen, Begründung, Auswirkungen.

## 8. Metadaten-Konventionen

Empfohlene minimale Frontmatter:

```yaml
---
title:
type:
status:
tags: []
created:
updated:
---
```

### Empfohlene `type`-Werte

* inbox
* daily
* project
* knowledge
* resource
* decision
* template

### Empfohlene `status`-Werte

* inbox
* draft
* active
* evergreen
* archived

## 9. Linking-Konzept

Wir nutzen drei Ebenen der Verknüpfung:

1. **Direkte Wiki-Links** auf konkrete Notizen
2. **Tags** für grobe thematische Bündelung
3. **Index-/Map-of-Content-Notizen** für kuratierte Einstiegspunkte

Regeln:

* lieber wenige sinnvolle Links als Link-Spam
* Knowledge Notes sollten auf Quellen und verwandte Konzepte verweisen
* Projekte sollten auf Entscheidungen, Meeting-Notizen und relevante Knowledge Notes zeigen

## 10. AI-Workflows mit Claude Code und GitHub Copilot CLI

## 10.0 Claude Code Skills

Das System definiert wiederholbare AI-Workflows als **Claude Code Skills** (Slash Commands), die über `pkm skill install` installiert werden.

Verfügbare Skills:

### `/capture-plan`

Erfasst den aktuellen Plan einer Claude Code Session als Inbox Note. Wird am Ende einer Coding Session aufgerufen, um die Planungsinhalte zu persistieren.

### `/distill-inbox`

Verarbeitet alle Inbox-Notizen nach dem Karpathy-Wiki-Muster:

1. `pkm process inbox --full` liefert alle Notizen + den aktuellen Stand von `04-knowledge/index.md`
2. Claude Code klassifiziert jede Notiz:
   - **Plan-Notizen** (`source: claude-code` + `# Plan:` Heading): werden an `02-projects/` weitergeleitet, nicht destilliert
   - **Reguläre Notizen**: werden als `distill`, `file-only` oder `skip` klassifiziert
3. Für jede Notiz wird ein Plan präsentiert — der Nutzer bestätigt vor der Ausführung
4. CLI-Kommandos (`pkm knowledge`, `pkm note move`, `pkm project update`) führen die Änderungen durch

## 10.0a Karpathy Wiki Pattern (Distillation)

Inspiriert durch Andrej Karpathys Ansatz für KI-gestützte Wissensbasen:

* Inbox-Notizen sind kurzlebige Inputs — sie landen in `00-inbox/` und werden verarbeitet
* Wertvolles Wissen wird als **Evergreen Prose** in `04-knowledge/<slug>.md` extrahiert — nicht als Zusammenfassung, sondern als zeitlose Erkenntnis
* `04-knowledge/index.md` ist die Topic-Übersicht — sie wächst mit jedem Distillations-Lauf
* `04-knowledge/log.md` ist der Audit-Trail — append-only, nie manuell bearbeiten
* Quellen werden nach der Destillation ins richtige Zielverzeichnis verschoben und archiviert

Anwendung: `pkm process inbox --full` + Claude Code (`/distill-inbox` Skill).

## 10.0b Projekt-Tracking

Plan-Notizen aus Claude Code Sessions (`source: claude-code` + `# Plan:` Heading) werden nicht destilliert, sondern an `02-projects/<slug>.md` weitergeleitet.

Jede Projekt-Notiz hat vier Abschnitte:
- **Intent** — dauerhaftes Ziel (stabil)
- **Current Status** — was zuletzt erarbeitet wurde (wird pro Session aktualisiert)
- **Next Steps** — was als nächstes zu tun ist (wird pro Session aktualisiert)
- **Plan History** — datumgestempelte Einträge mit Wikilinks zu archivierten Plan-Notizen (append-only)

Anwendung: `pkm project update <slug>` + `pkm project list`.

## 10.1 ChatGPT Export -> Inbox Note

Ziel: wertvolle Chat-Inhalte in saubere Markdown-Notizen überführen.

Der Workflow erzeugt zunächst immer eine Inbox-Note.

Der Workflow soll:

* eine Chat-Antwort oder einen Chat-Ausschnitt als Eingabe nehmen
* Kernaussagen extrahieren
* irrelevante Gesprächsteile entfernen
* einen vorläufigen Titel erzeugen
* minimal sinnvolle Struktur erzeugen
* eine importfähige Markdown-Datei in `00-inbox/` erzeugen

Die Klassifikation in Knowledge-, Decision- oder Project-Notes erfolgt erst im Inbox-Processing-Schritt.

## 10.2 Coding-Session Export -> Inbox Note

Ziel: Ergebnisse aus Claude Code oder GitHub Copilot CLI dauerhaft dokumentieren.

Der Workflow erzeugt zunächst immer eine Inbox-Note.

Der Workflow soll:

* eine Session, Konsole, Diff-Beschreibung oder Zusammenfassung als Input nutzen
* Problem, Ursache, Lösung und Auswirkungen strukturieren
* Learnings extrahieren
* offene Follow-ups markieren
* eine Markdown-Datei in `00-inbox/` erzeugen

Eine spätere Klassifikation als Troubleshooting-, Decision-, Architecture- oder Knowledge-Note erfolgt im Inbox-Processing.

## 10.3 Inbox Processing (zentraler AI-Schritt)

Ziel: rohe Notizen aus `00-inbox/` in strukturierte Wissensartefakte überführen.

Dies ist der wichtigste AI-Workflow im gesamten System.

Claude Code oder GitHub Copilot CLI sollen dabei helfen:

* Notiztyp zu erkennen
* Titel zu verbessern
* fehlendes Frontmatter ergänzen
* Tags vorzuschlagen
* interne Links vorzuschlagen
* relevante bestehende Notizen zu referenzieren
* zusätzliche Struktur zu ergänzen
* offene Fragen zu extrahieren
* To-dos abzuleiten
* Zielordner vorzuschlagen
* optional die Datei automatisch zu verschieben (nach Bestätigung)

Dieser Schritt trennt bewusst:
Capture (schnell)
von
Curation (qualitätsgesichert).

## 10.4 Note Distillation

Ziel: aus umfangreichen Notizen eine wertvolle Evergreen Note machen.

Claude Code oder GitHub Copilot CLI sollen:

* Kernaussagen extrahieren
* Redundanzen entfernen
* klare Struktur erzeugen
* eine prägnante Zusammenfassung schreiben
* mögliche Verbindungen zu anderen Themen hervorheben

## 10.4a Readwise Reader Sync

Gespeicherte Artikel aus Readwise Reader werden inkrementell in `00-inbox/` synchronisiert:

```bash
pkm sync readwise           # inkrementeller Sync (State-File)
pkm sync readwise --dry-run
pkm sync readwise auth      # Token einrichten
```

Jeder Artikel wird als Resource Note mit Highlights angelegt. Der Sync ist idempotent (URL-basierte Duplikaterkennung) und speichert den letzten Sync-Zeitpunkt in `~/.config/pkm/readwise_sync_state.json`.

## 10.5 Project Memory Support

Claude Code oder GitHub Copilot CLI sollen für Projektordner:

* Meeting-Notizen konsolidieren
* offene Punkte sammeln
* Entscheidungen extrahieren
* Status-Updates formulieren
* nächste Schritte ableiten

## 10.6 Standard Note Creation

Claude Code soll aus kurzen Eingaben oder Stichpunkten standardisierte Notizen erzeugen.

Beispiele:

* Erzeuge die Daily Note für heute
* Formatiere diese Stichpunkte als Meeting Note
* Leite aus diesem Verlauf eine Decision Note ab

## 10.7 Knowledge Refactoring

Claude Code oder GitHub Copilot CLI sollen helfen,

* große Notizen zu splitten
* doppelte Inhalte zu erkennen
* schwache Titel zu verbessern
* Querverweise zu ergänzen
* veraltete Inhalte zu markieren

## 11. Beispiel-Prompts für Claude Code und GitHub Copilot CLI

### ChatGPT-Zusammenfassung exportieren

> Verdichte den folgenden Chat-Verlauf zu einer sauberen Markdown-Notiz für Obsidian. Entferne Small Talk und Redundanzen. Erzeuge einen klaren Titel, einen passenden Notiztyp, Frontmatter, Kernaussagen, offene Fragen und empfohlene interne Links.

### Coding-Session dokumentieren

> Erstelle aus dieser Coding-Session eine Markdown-Notiz für Obsidian. Strukturiere in Problem, Kontext, Root Cause, Lösung, betroffene Komponenten, Learnings und nächste Schritte. Wenn eine Architekturentscheidung enthalten ist, ergänze einen Abschnitt Entscheidung.

### Inbox aufräumen

> Analysiere die Dateien in `00-inbox/`. Schlage für jede Notiz einen passenden Titel, Notiztyp, Zielordner, Tags und relevante interne Links vor. Erstelle keine Änderungen ohne sie explizit aufzulisten.

### Knowledge Note erstellen

> Verdichte diese Rohnotiz zu einer prägnanten Knowledge Note. Behalte nur wiederverwendbare Erkenntnisse, strukturiere in kurze Abschnitte und ergänze am Ende offene Fragen.

### Projektstatus extrahieren

> Analysiere alle Markdown-Dateien in diesem Projektordner und erstelle eine Statusübersicht mit Ziel, aktuellem Stand, Risiken, Entscheidungen und nächsten Schritten.

### Decision Record ableiten

> Prüfe diese Notizen auf explizite oder implizite Entscheidungen und formuliere daraus einen Decision Record im Template-Format.

### Daily Note erzeugen

> Erstelle die Daily Note für heute auf Basis dieses Templates. Übernimm die folgenden Rohnotizen, gruppiere sie in sinnvolle Abschnitte und extrahiere Aufgaben, Follow-ups und offene Fragen.

### Meeting Note erzeugen

> Formatiere diese Stichpunkte in eine strukturierte Meeting Note mit Teilnehmern, Kontext, besprochenen Punkten, Entscheidungen, Risiken und nächsten Schritten.

## 12. Templates

Es sollen mindestens folgende Templates existieren:

* Inbox Note
* Daily Note
* Meeting Note
* Decision Note
* Knowledge Note
* Troubleshooting Note
* Project Note
* Resource Note

Die Templates sollen bewusst einfach, markdown-first und AI-friendly sein.

## 12.1 Inbox Note Template

Zweck:
Einheitliches Intake-Format für manuelle Captures sowie Exporte aus ChatGPT, Claude Code und GitHub Copilot CLI.

```markdown
---
title: 
type: inbox
status: inbox
source: manual
created: 
updated: 
tags: []
type_hint: 
---

## Context

- Herkunft:
- Anlass:
- Erwarteter Zieltyp:

## Content


## Open Questions

- 

## Follow-ups

- 
```

Hinweise:

* `type` bleibt zunächst immer `inbox`
* `type_hint` ist optional und dient nur als Vorab-Hinweis
* Die eigentliche Klassifikation erfolgt erst im Inbox-Processing

## 12.2 Daily Note Template

Zweck:
Täglicher Einstiegspunkt für Gedanken, Arbeit, Learnings und Follow-ups.

Dateiname:
`YYYY-MM-DD.md`

```markdown
---
title: 
type: daily
status: active
created: 
updated: 
tags: [daily]
---

# Daily Note

## Focus Today

- 

## Schedule / Meetings

- 

## Notes

- 

## Decisions / Learnings

- 

## Open Loops

- 

## Tasks

- [ ] 
```

## 12.3 Meeting Note Template

Zweck:
Strukturierte Erfassung von Meetings, Syncs, Workshops und Abstimmungen.

```markdown
---
title: 
type: meeting
status: draft
created: 
updated: 
tags: [meeting]
date: 
participants: []
project: 
source: manual
---

# Meeting Note

## Context

- Date:
- Participants:
- Project / Area:
- Purpose:

## Agenda / Topics

- 

## Notes

- 

## Decisions

- 

## Risks / Issues

- 

## Next Steps

- [ ] 

## Related Notes

- 
```

Hinweise:

* Meeting Notes können manuell oder per `pkm meeting create` erzeugt werden
* Standardmäßig können sie zunächst in der Inbox landen, wenn sie noch kuratiert werden sollen

## 12.4 Decision Note Template

Zweck:
Dokumentation von Architektur-, Tooling-, Prozess- oder Arbeitsentscheidungen.

```markdown
---
title: 
type: decision
status: draft
created: 
updated: 
tags: [decision]
decision_date: 
project: 
related_notes: []
---

# Decision

## Status

Draft / Accepted / Superseded

## Context


## Decision


## Options Considered

- Option A:
- Option B:

## Rationale


## Consequences

### Positive

- 

### Negative / Trade-offs

- 

## Follow-ups

- [ ] 

## References

- 
```

Hinweise:

* gut geeignet für aus Meetings oder Coding-Sessions abgeleitete Entscheidungen
* kann zunächst in der Inbox entstehen und später in `06-decisions/` verschoben werden

## 12.5 Knowledge Note Template

Zweck:
Verdichtete, wiederverwendbare Erkenntnis zu einem Thema.

```markdown
---
title: 
type: knowledge
status: evergreen
created: 
updated: 
tags: []
related_notes: []
---

# Summary


## Key Points

- 

## Explanation


## Examples / Applications

- 

## Related Concepts

- 

## Open Questions

- 

## Sources

- 
```

Hinweise:

* Fokus auf Wiederverwendbarkeit statt Rohdokumentation
* ideal als Zielstruktur nach Distillation aus Inbox Notes

## 12.6 Troubleshooting Note Template

Zweck:
Dokumentation von Fehlerbildern, Ursachenanalyse und Lösungsschritten.

```markdown
---
title: 
type: troubleshooting
status: draft
created: 
updated: 
tags: [troubleshooting]
systems: []
project: 
source: claude-code
related_notes: []
---

# Troubleshooting

## Problem


## Context

- Environment:
- System / Component:
- Trigger:

## Symptoms

- 

## Root Cause


## Resolution

1. 
2. 
3. 

## Validation

- 

## Learnings

- 

## Follow-ups

- [ ] 

## References

- 
```

Hinweise:

* besonders nützlich für Claude-Code- oder Copilot-CLI-Exporte
* kann später in `04-knowledge/`, `02-projects/` oder einen dedizierten Troubleshooting-Bereich überführt werden

## 12.7 Project Note Template

Zweck:
Persistente Projektverfolgung mit schnellem Kontext-Switch zwischen Projekten. Wird über `pkm project update` gepflegt, nicht manuell erstellt.

```markdown
---
title: 
type: project
status: active
source: claude-code
created: 
updated: 
tags: [project]
---

## Intent

Was das Projekt erreichen soll (stabil, selten geändert).

## Current Status

Was zuletzt erarbeitet wurde (wird pro Session aktualisiert).

## Next Steps

- [ ] Nächste Aktion
- [ ] Weitere Aktion

## Plan History

### YYYY-MM-DD — Session-Beschreibung

[[YYYY-MM-DD-HHMM-claude-code-plan-notiz]]
```

Die vier Abschnitte:
- **Intent** — dauerhaftes Projektziel
- **Current Status** — Zustand nach der letzten Session
- **Next Steps** — Checkliste für die nächste Session
- **Plan History** — append-only Verlauf mit Wikilinks zu archivierten Plan-Notizen

## 12.8 Resource Note Template

Zweck:
Dokumentation externer Quellen wie Artikel, Videos, Docs, Bücher oder Links.

```markdown
---
title: 
type: resource
status: draft
created: 
updated: 
tags: [resource]
url: 
author: 
published: 
related_notes: []
---

# Resource

## Source

- URL:
- Author:
- Published:

## Summary


## Why It Matters


## Key Extracts

- 

## Related Notes

- 
```

## 12.9 Template-Prinzipien

Die Templates folgen folgenden Regeln:

1. so wenig Pflichtfelder wie möglich
2. klare Abschnittsüberschriften für AI-gestützte Verarbeitung
3. markdown-only, keine versteckten Obsidian-spezifischen Mechanismen als Voraussetzung
4. kompatibel mit manuellem Schreiben und CLI-Generierung
5. geeignet für spätere Anreicherung durch `pkm process inbox`

## 13. MVP (implementiert)

Der MVP ist vollständig umgesetzt. Stand 2026-04:

* Obsidian Vault mit definierter Ordnerstruktur und Templates
* `pkm capture` — schnelles Erfassen in `00-inbox/`
* `pkm process inbox --full` — Kontext-Ausgabe für Claude Code Distillation Sessions
* `pkm knowledge` — Karpathy-Wiki-Distillation: append-topic, update-index, append-log
* `pkm project` — Projektverfolgung mit Intent / Status / Next Steps / Plan History
* `pkm note move` — Verschieben mit automatischer Frontmatter-Aktualisierung
* `pkm daily / meeting / decision create` — strukturierte Notizen aus Templates
* `pkm sync readwise` — inkrementeller Readwise Reader Sync mit Highlights
* Claude Code Skills: `/capture-plan`, `/distill-inbox`
* Erste Distillation durchgeführt: 49 Notizen → 8 Knowledge Topic Pages

## 14. Spätere Ausbaustufen

### Phase 2

* automatische Link-Vorschläge
* Qualitätschecks für Frontmatter und Dateinamen
* semiautomatische Refactoring-Skripte
* Projektzusammenfassungen per Befehl

### Phase 3

* lokale Suche mit Embeddings / RAG über ausgewählte Inhalte
* Retrieval-gestützte Prompt-Workflows
* AI-gestützte Wissenslandkarten
* Integration von E-Mail, Web-Clips oder PDFs

## 15. Offene Entscheidungen

1. Soll die Struktur eher PARA-orientiert oder stärker wissensorientiert sein?
2. Wie viel YAML-Frontmatter ist sinnvoll, ohne die Pflege zu schwer zu machen?
3. Sollen Tasks direkt in Obsidian gepflegt werden oder in einem externen System?
4. Soll Git von Anfang an verpflichtend sein?
5. Welche Obsidian-Plugins sind für den MVP erlaubt?
6. Wie strikt sollen Claude Code und GitHub Copilot CLI Änderungen automatisch anwenden dürfen?
7. Exporte aus ChatGPT und Coding-Agents gehen immer zunächst in die Inbox (bereits festgelegt)
8. Brauchen wir separate Templates für Troubleshooting Notes und Architecture Notes?

## 16. Erfolgsmetriken

Das System ist erfolgreich, wenn:

* neue Informationen schnell erfasst werden können
* Wissen später zuverlässig wiedergefunden wird
* AI die Pflege spürbar beschleunigt
* Notizen über Zeit eher besser statt chaotischer werden
* die Struktur auch nach mehreren Monaten noch verständlich bleibt

## 17. Nächste Schritte

MVP ist abgeschlossen. Nächste Ausbaustufen:

1. Unit-Tests für `internal/readwise/` schreiben
2. Raycast-Integration für schnelles Capture ohne Terminal
3. `pkm note refine` — AI-gestützte Überarbeitung bestehender Notizen
4. `pkm note link` — Vorschläge für interne Wikilinks
5. Automatische Tag-Vergabe beim Capture oder Inbox-Processing
6. Retrieval-gestützte Workflows (lokale Suche, Embeddings)


## 18. Geplant: GitHub, Markdown und MCP

**Status:** Entwurf / noch nicht implementiert. **Erstellt und aktualisiert:** 2026-10-11.
Grundlage ist die Unterhaltung „Daily AI Tech Brief“ vom 2026-10-11 zum
Reader→PKM-Workflow und zur Vereinheitlichung mit Live Journal und Training Journal.
Dieser Abschnitt spezifiziert das Zielverhalten; er behauptet keine vorhandenen MCP-Funktionen.

### 18.1 Ziel, Zuständigkeiten und Abgrenzung

Das PKM wird ein langlebiger **Knowledge Layer** für eigene, verdichtete und
wiederverwendbare Erkenntnisse. **Readwise Reader bleibt der Source Layer** für
Leseliste, Originalmaterial, Lesestatus, Highlights und artikelbezogene Notizen.
Ein archivierter Reader-Artikel ist weiterhin eine auffindbare Quelle; Archivierung
ist weder Löschung noch ein Signal für automatische Übernahme ins PKM.

| System | Primäre Struktur | Verantwortlich für |
|---|---|---|
| Readwise Reader | Dokument und Leseprozess | Sammeln, Lesen, Markieren, Archivieren |
| Live Journal | Zeit und Ereignisse | Persönlicher Verlauf und Reflexion |
| Training Journal | Trainingseinheit und Zeit | Trainingsverlauf, Beobachtungen und Auswertung |
| PKM | Aussage, Konzept und Verknüpfung | Übertragbare Erkenntnisse, Belege, Anwendung und Themenübersichten |
| Projekt-Repository | Projekt und Entscheidung | Specs, ADRs, Aufgaben und projektspezifischer Arbeitsstand |

GitHub + Markdown + MCP ist das gemeinsame Persistenz- und Zugriffsmuster, kein
Grund für ein gemeinsames Datenschema oder Repository aller Journale. Eine
Journalbeobachtung darf eine PKM-Erkenntnis begründen und wird dann verlinkt.
Tagesverlauf, Trainingsprotokolle und Projektpläne werden nicht vollständig ins PKM
kopiert. Eine Erkenntnis kann ohne externe Quelle aus eigener Erfahrung entstehen;
ihr Erfahrungskontext muss dann benannt sein.

Nicht Bestandteil der ersten Version sind autonome Wissenspflege, vollständiger
Reader-Import, bidirektionale Reader-Synchronisierung, Journalfunktionen,
Aufgabenverwaltung, Vektordatenbank oder ein eigener RAG-Stack.

### 18.2 Anforderungen

- **R1 – Persistenz:** Versionierte Markdown-Dateien mit YAML-Frontmatter in einem
  privaten GitHub-Wissensrepository sind die Source of Truth. Indizes und Caches
  müssen daraus vollständig rekonstruierbar sein.
- **R2 – Zugriff:** MCP ist die primäre API auf das Wissen. Clients können suchen,
  lesen, erfassen, erstellen, aktualisieren und verknüpfen. Lesen und manuelles
  Bearbeiten mit einem beliebigen Markdown-Editor bleiben möglich.
- **R3 – Kuratierung:** Nur bewusst ausgewählte Quellen und wiederverwendbare
  Erkenntnisse werden übernommen. Capture und fachliche Freigabe sind getrennt.
- **R4 – Nachvollziehbarkeit:** Aussagen, Quellenzitate und eigene Interpretation
  bleiben unterscheidbar. Herkunft, interne Verweise und Git-Version sind abrufbar.
- **R5 – Integrität:** Validierung, stabile IDs, idempotente Wiederholungen und
  Konflikterkennung verhindern Duplikate, ungültige Links und verlorene Änderungen.
- **R6 – Gemeinsamer Kern:** MCP und eine gegebenenfalls weitergeführte CLI verwenden
  dieselben Regeln für Storage, Validierung, Suche und Linking.
- **R7 – Migration:** Bestehende Notizen und ihre Herkunft bleiben erhalten. Ein
  Schema- oder Strukturwechsel erfolgt mit Vorschau, Prüfung und Rückweg.

### 18.3 Architektur und Repository-Struktur

```text
Readwise Reader / eigene Gedanken / ausgewählte Journalbeobachtungen
                       ↓ bewusste Auswahl und Verdichtung
ChatGPT / Coding Agents / optionale CLI
                       ↓ MCP
                    pkm-mcp
                       ↓ gemeinsamer pkm-core
              privates GitHub-Wissensrepo
                       ↓
              Markdown + YAML-Frontmatter
```

`pkm.ai` ist heute das **Software-Repository der Go-CLI**, nicht der Wissensspeicher.
Diese Spec bleibt hier. Der geplante Wissensspeicher ist ein separates privates
Repository, Arbeitsname `pkm`. Ein separates Software-Repository `pkm-mcp` ist
vorgesehen; endgültige Namen, Deployment und die physische Paketaufteilung des
`pkm-core` werden vor der Implementierung festgelegt. Diese Spezifikation erstellt
weder Repositories noch einen MCP-Dienst und verschiebt keine bestehenden Daten.

Geplante Struktur des Wissensrepositorys:

```text
README.md                 # Zweck, Schema, Arbeitsweise und Einstiegspunkte
inbox/                    # noch zu prüfende Notizen aller drei Typen
notes/                    # kuratierte insight-Notizen
sources/                  # ausgewählte source-Notizen
topics/                   # kuratierte Themenübersichten
templates/
  insight.md
  source.md
  topic.md
```

Die Ordner bleiben zunächst flach. Fachliche Zuordnung erfolgt über Topics und
Links; zusätzliche Domänenordner werden erst bei tatsächlichem Bedarf eingeführt.
Jede Datei heißt `<id>-<english-ascii-kebab-slug>.md`. Die unveränderliche ID ist eine
UUID; Titel, Slug und Pfad dürfen sich ändern. `inbox/` ist eine Bearbeitungsstufe,
kein vierter Note-Typ. Freigegebene Notizen wechseln in den typabhängigen Zielordner.
Archivierte Notizen bleiben dort auffindbar und erhalten `status: archived`.

Obsidian ist optionaler Editor und Browser. Das Schema und die Nutzung über MCP
setzen weder Obsidian noch dessen Plugins voraus.

### 18.4 Datenmodell und Frontmatter

| Typ | Inhalt | Empfohlene Abschnitte |
|---|---|---|
| `insight` | Eine eigenständig verständliche Aussage oder ein Prinzip | These, Kontext/Geltungsbereich, Belege, Konsequenzen/Anwendung, Grenzen/offene Fragen |
| `source` | Schlanker Nachweis einer bewusst ausgewählten Quelle | Quelle, Relevanz, ausgewählte Highlights, eigene Gedanken |
| `topic` | Kuratierter Einstieg in ein Thema mit Links zu Insights | Überblick, Leitfragen, Erkenntnisse, Widersprüche/offene Fragen |

Topics ersetzen keine einzelnen Insights. Eine Source ist kein Artikelarchiv;
nur relevante Auszüge und Kontext werden gespeichert. Ein Insight kann mehrere
Sources nutzen, eine Source mehrere Insights belegen. Themen sind ebenfalls
Notizen mit stabilen IDs und können mehrere Insights bündeln.

Gemeinsame Pflichtfelder:

| Feld | Typ und Bedeutung |
|---|---|
| `schema_version` | Integer, anfangs `1`; unbekannte Versionen nicht schreibend verarbeiten |
| `id` | Eindeutige, unveränderliche UUID |
| `title` | Nicht leerer String |
| `type` | `insight`, `source` oder `topic` |
| `status` | `inbox`, `draft`, `active` oder `archived` |
| `created`, `updated` | Datum `YYYY-MM-DD`; `created` unveränderlich, `updated` bei inhaltlicher Änderung |
| `origin` | `manual`, `readwise`, `chat`, `journal` oder `migration` |
| `topics` | Liste von IDs existierender `topic`-Notizen; leer zulässig |
| `sources` | Liste von IDs existierender `source`-Notizen; leer bei eigener Erkenntnis zulässig |
| `related` | Liste anderer vorhandener Note-IDs; leer zulässig |
| `tags` | Freie Schlagwörter; Liste, leer zulässig |

Source-spezifisch: `url` ist der dauerhafte Herkunftslink; bei nicht webbasierten
Quellen ist stattdessen `citation` als bibliografischer Nachweis erforderlich.
Optional sind `author`, `published`, `reader_document_id` und `reader_url`.
Reader-Captures speichern die Dokument-ID und den Reader-Link, soweit verfügbar,
sowie die Original-URL. Unbekannte Metadaten werden weggelassen, nicht erfunden.
`origin` beschreibt den Erfassungsweg; `sources` verweist auf Belege.

Optionale `legacy_path`- und `legacy_type`-Felder unterstützen die Migration.
Erweiterungsfelder bleiben beim Aktualisieren erhalten. Das alte skalare Feld
`source: readwise` wird bei der Migration zu `origin: readwise`, nicht zu `sources`.

Beispiel einer noch unbelegten eigenen These, bewusst als Entwurf markiert:

```yaml
---
schema_version: 1
id: "9d7b83c8-57e3-4aac-b180-81557d017544"
title: "Agent-Infrastruktur als Plattformfunktion"
type: insight
status: inbox
created: "2026-10-11"
updated: "2026-10-11"
origin: manual
topics: []
sources: []
related: []
tags: [ai-agents, platform-engineering]
---
```

Der Body formuliert etwa die zu prüfende These, dass Identity, Gateways und Policy
als Plattformdienste konsumiert werden können, und beschreibt Einsatzgrenzen.
Vor Aktivierung werden belastbare Belege oder der eigene Erfahrungskontext ergänzt.
Eine Produktankündigung allein ist noch kein allgemeingültiger Architekturbeleg.

Lebenszyklus: Capture erzeugt `inbox`; die Bearbeitung kann zu `draft` wechseln.
Fachliche Prüfung und bewusste Freigabe führen zu `active`. `archived` kennzeichnet
nicht mehr gepflegte Inhalte; Reaktivierung ist möglich. Eine neu erzeugte
AI-Synthese wird nicht automatisch aktiv. Der Status wird nicht aus dem Ordnernamen
oder dem Reader-Lesestatus abgeleitet. `inbox` liegt in `inbox/`, die übrigen Status
liegen im jeweiligen Typordner; Statuswechsel und Verschieben sind atomar.

Metadaten enthalten ID-Verweise als verbindliche Beziehungen. Im Body werden
zusätzlich relative Markdown-Links mit lesbaren Titeln verwendet, die auf GitHub
und lokal funktionieren. Pfadänderungen müssen betroffene interne Links mitführen;
IDs bleiben stabil. Backlinks werden aus den ausgehenden Beziehungen abgeleitet,
nicht als zweite manuell zu pflegende Wahrheit gespeichert.

### 18.5 Reader→PKM-Workflow

1. **Lesen:** Links landen in Reader. Beim Lesen werden wenige relevante Stellen
   markiert und eigene Gedanken ergänzt. Die meisten Artikel gehen anschließend
   ausschließlich ins Reader-Archiv.
2. **Auswählen:** Ein PKM-Kandidat enthält etwas, das später für Entscheidungen,
   Präsentationen, Projekte oder eigenes Verständnis wiederverwendbar ist.
   „Interessant“ oder „neu veröffentlicht“ allein reicht nicht.
3. **Abgleichen:** Vor Neuanlage nach vorhandenen Sources, Insights und Topics
   suchen. Ähnliche Aussagen bevorzugt ergänzen, Gegenbelege sichtbar aufnehmen.
4. **Quelle erfassen:** `capture_source` übernimmt Original-Link, Reader-Referenz,
   ausgewählte Highlights und eigene Notizen als schlanke Source in die Inbox.
   Zitate werden mit Herkunft kenntlich gemacht; eigene Gedanken stehen getrennt.
5. **Verdichten:** Aus einer oder mehreren Quellen eine selbstständige These mit
   Geltungsbereich, Belegen und Konsequenzen formulieren oder einen bestehenden
   Insight aktualisieren. Mehrere Artikel dürfen eine gemeinsame Erkenntnis ergeben.
6. **Prüfen und verknüpfen:** Entwurf lesen, Aussagen gegen Quellen prüfen, passende
   Topics und verwandte Notes verlinken. Die bestätigte Fassung aktivieren.
7. **Abschließen:** Erfolgreiche Speicherung liefert Note-IDs und Git-Revision.
   Reader-Archivierung bleibt ein unabhängiger Leseschritt und wird vom PKM nicht
   automatisch ausgelöst. Ein bereits archivierter Artikel kann später noch
   bewusst als Quelle erfasst werden.

Es gibt in Version 1 keinen automatischen Import der Reader-Inbox und keinen
Import allein aufgrund eines Highlights. Die Reader-Auswahl kann über den
bestehenden Reader-Zugang des Clients erfolgen; `pkm-mcp` muss zum Capture nicht
selbst auf Reader zugreifen. Ein fehlender Reader-Zugriff blockiert das Erfassen
bekannter URLs und ausdrücklich übergebener Auszüge nicht.

### 18.6 MCP-Funktionen und Verhaltensvertrag

Die folgenden Namen beschreiben den geplanten Tool-Vertrag, keine bereits
verfügbare Schnittstelle. Generative Verdichtung erfolgt zunächst im aufrufenden
AI-Client; der MCP validiert und persistiert die übergebenen Inhalte.

| Tool | Wesentliche Eingaben | Ergebnis und Verhalten |
|---|---|---|
| `search_notes` | `query`, Filter für Typ/Status/Topic/Tags, `limit`, `cursor` | Volltextsuche über Titel, Body und Metadaten; Treffer mit ID, Titel, Pfad, Auszug, Status und Revision; paginiert, archivierte Notes nur auf Wunsch |
| `get_note` | `id` oder relativer `path` | Vollständiges Frontmatter und Markdown, aufgelöste Beziehungen und Revision; eindeutiger Fehler bei fehlender Note |
| `list_topics` | optional Suchtext, `limit`, `cursor` | Vorhandene Topics mit ID, Titel und Kurzbeschreibung zur Wiederverwendung |
| `create_note` | Typ, Titel, Inhalt, Metadaten, `request_id`, `dry_run` | Validierte neue Note mit generierter ID, initial `inbox`; keine stille Überschreibung |
| `update_note` | ID, explizite Änderungen an Body/Metadaten, `expected_revision`, `request_id`, `dry_run` | Nur angegebene Inhalte ändern; unbekannte Metadaten erhalten; Statuswechsel validieren |
| `link_notes` | Ausgangs-ID, Ziel-ID, Beziehung `related`/`sources`/`topics`, Aktion `add`/`remove`, `expected_revision`, `request_id`, `dry_run` | Typgerechte Beziehung ohne Duplikate; fehlende Ziele und Selbstverweise ablehnen |
| `capture_source` | URL oder bibliografischer Nachweis, Titel, optionale Reader-ID/-URL, ausgewählte Highlights und eigene Notizen, `request_id`, `dry_run` | Source anlegen oder vorhandene Source zurückgeben; bei geändertem Material Vorschau und konfliktgeprüftes Update statt blindem Überschreiben |
| `create_insight_from_sources` | Source-IDs, vom Client formulierter Entwurf, Titel, Topics, `request_id`, `dry_run` | Quellen prüfen und einen verknüpften Insight als `inbox` erstellen; keine erfundenen Belege oder automatische Aktivierung |

Alle Schreibtools unterstützen eine Vorschau (`dry_run`) ohne Dateischreibvorgang,
Commit oder Push. Bereits erteilte, konkrete Schreibaufträge dürfen direkt
umgesetzt werden; zusätzliche AI-Vorschläge zur Umstrukturierung müssen zuvor
prüfbar sein. Die bestehende CLI-Regel zur Bestätigung bei Inbox-Verarbeitung
bleibt bis zu ihrer ausdrücklich implementierten Ablösung bestehen.

`expected_revision` bezeichnet den vom letzten Lesezugriff bekannten Repository-
Commit. Stimmt dieser nicht mehr mit dem aktuellen Stand überein, erfolgt ein
Konfliktfehler ohne Schreibvorgang. Der Client liest neu und überprüft seinen Diff.
Ein konservativer Konflikt bei Änderungen an anderen Dateien ist in Version 1 zulässig.

Erfolgreiche Schreibantworten enthalten ID, Pfad, Änderungszusammenfassung und
Commit-SHA; No-ops werden ausdrücklich ausgewiesen. Fehler unterscheiden mindestens
Validierungsfehler, unbekannte ID, Konflikt, fehlende Berechtigung und Storage-/Sync-
Fehler. Fehlgeschlagene Operationen erzeugen keine teilweise aktive Note.

`request_id` identifiziert eine logische Mutation: gleiche ID und gleiche Nutzlast
liefern beim Retry dasselbe Ergebnis, eine andere Nutzlast unter derselben ID einen
Konflikt. Der Idempotenznachweis muss auch nach Neustart rekonstruierbar sein, etwa
aus Commit-Metadaten; ein flüchtiger Cache allein reicht nicht.

Source-Deduplizierung verwendet zuerst eine vorhandene Reader-Dokument-ID, sonst
die konservativ normalisierte Original-URL. Nur bekannte Tracking-Parameter werden
entfernt; inhaltlich relevante URL-Bestandteile bleiben erhalten. Bibliografische
Quellen ohne URL werden bei unklarer Identität zur manuellen Prüfung vorgeschlagen.
Gleiche Titel allein begründen keine Zusammenführung. Bereits kuratierte eigene
Gedanken und Highlights werden bei erneutem Capture niemals still ersetzt.

Die Suche startet mit nachvollziehbarer lexikalischer Suche. Ein Index ist ein
abgeleiteter Cache und an einen Commit gebunden. Nach erfolgreicher Mutation muss
die anschließende Suche diesen Stand sehen; veraltete Indizes werden aktualisiert
oder die Dateien direkt gelesen. Ein Suchtreffer darf nicht als vollständiger
Kontext für eine Überschreibung verwendet werden.

### 18.7 Versionierung, Git-Workflow und Zugriff

- Das private Wissensrepository ist die kanonische Historie. Der MCP arbeitet auf
  einem konfigurierten Branch; Standard ist der Default-Branch. GitHub-API oder
  lokaler Clone sind mögliche Adapter mit demselben Konsistenzvertrag.
- Vor dem Schreiben wird der aktuelle Remote-Stand gelesen. Validierung und Diff
  gehen dem Commit voraus. Eine fachlich zusammengehörige Operation einschließlich
  Link- und Pfadkorrekturen erzeugt genau einen atomaren Commit.
- Der Commit nennt Aktion und betroffene Note-IDs; Request-ID und Client-Herkunft
  ermöglichen Zuordnung. Die konfigurierte Git-Identität wird verwendet, keine
  erfundene Autorenidentität. Fachliche Provenienz steht in den Notes.
- Remote-Aktualisierung erfolgt nur auf Basis der erwarteten Revision. Keine
  Force-Pushes, kein stilles Last-write-wins. Bei Push-Fehlern wird lokale Speicherung
  nicht als erfolgreicher GitHub-Abschluss ausgegeben; Retry prüft, ob die Operation
  möglicherweise bereits angekommen ist.
- Explizite einzelne Capture-/Pflegeaufträge können direkt auf dem Arbeitsbranch
  landen. Massenmigration und größere Umstrukturierung erfolgen in einem separaten
  Branch mit geprüftem Diff und Pull Request vor Übernahme in den Default-Branch.
- Direkte manuelle Markdown-Edits bleiben erlaubt und müssen dieselben Schema- und
  Linkprüfungen bestehen. Fremde uncommittete Änderungen werden nicht mitcommittet
  oder überschrieben. Der MCP benötigt dafür einen isolierten Arbeitsbereich.
- Rücknahmen erfolgen per Revert, ohne Historie umzuschreiben. Archivieren wird
  Löschen vorgezogen; ein Löschtool gehört nicht zur ersten Version.
- GitHub-Zugangsdaten liegen außerhalb des Wissensrepositorys. Zugriff ist auf das
  konfigurierte Repo beschränkt; Pfade außerhalb der erlaubten Notizbereiche,
  Path Traversal und Symlink-Ausbrüche werden abgelehnt. Lesen und Schreiben
  benötigen entsprechend begrenzte Berechtigungen.
- Quellen und Notizinhalte sind Daten, keine Ausführungsanweisungen. Capture führt
  weder enthaltene Befehle aus noch folgt es automatisch beliebigen Links.

### 18.8 Bestehende CLI und Übergang

Die CLI wird nicht sofort entfernt. Wiederverwendbare Logik wird in einen gemeinsamen
Kern für Storage, Validierung, Indexierung und Linking überführt. `pkm-mcp` nutzt
diesen Kern; eine weitergeführte CLI verwendet denselben Kern oder wird MCP-Client.
Es entstehen keine unabhängig gepflegten Geschäftsregeln pro Zugang.

Die aktuelle CLI erwartet nummerierte Vault-Ordner und andere Note-Typen. Ihr
Frontmatter-Modell unterstützt die neuen Felder noch nicht vollständig. Sie darf
daher nicht ungeprüft auf den neuen Wissensspeicher schreiben: Alte Befehle könnten
Typen umdeuten oder Metadaten verlieren. Bis zur Anpassung bleibt der alte Vault
getrennt und die bestehende CLI dort unverändert nutzbar.

| Bisher | Geplanter Umgang |
|---|---|
| `00-inbox/` | Bewusst ausgewählte Inhalte als typisierte Inbox-Entwürfe übernehmen |
| `04-knowledge/<slug>.md` | Inhalt prüfen: einzelne Aussagen als `insight`, Themenübersichten als `topic`; keine pauschale Umbenennung |
| `04-knowledge/index.md` | Als Ausgangspunkt für Topics und deren Verknüpfungen verwenden |
| `04-knowledge/log.md` | Alten Audit-Trail unverändert bewahren; neue Änderungen sind über Git nachvollziehbar |
| `05-resources/` | Nur weiterhin relevante Quellen als schlanke `source` übernehmen; Herkunft und eigene Gedanken erhalten |
| `01-daily/`, Meetings, Training | Im bisherigen Bestand bewahren; Journalzuständigkeit getrennt klären; nur Erkenntnisse ins PKM extrahieren |
| Projekte, Areas, Decisions | Bestehenden Kontext bewahren; Projekt-Specs/ADRs bleiben im Projekt, nur übertragbares Wissen ins PKM |
| `07-templates/` | Nach Schemafreigabe drei neue Templates ergänzen bzw. ersetzen |
| `09-archive/`, Attachments | Kein pauschaler Import oder Löschen; relevante Verweise vor Migration prüfen |
| `pkm sync readwise` | Bestehenden Inbox-Sync zunächst unverändert lassen; im neuen Workflow durch explizites Source-Capture ersetzen |

Migration: Bestand und Links inventarisieren, Ausgangsrevision sichern, eine
repräsentative kleine Auswahl als Vorschau konvertieren, Zuordnung alter Pfade zu
neuen IDs dokumentieren, Metadaten und Links validieren, danach schrittweise per
Review übernehmen. Der alte Vault bleibt bis zur geprüften Vollständigkeit erhalten;
kein gleichzeitiger Schreibbetrieb alter CLI und neuer API auf denselben Dateien.

Die CLI kann entfallen, wenn Capture, Suche, Lesen, Pflege und Source-Übernahme über
MCP im Alltag funktionieren und keine benötigten lokalen oder Offline-Abläufe mehr
von ihr abhängen. Andernfalls bleibt sie ein schlanker optionaler Zugang. Entfernung,
Kompatibilitätsdauer und Ersatz bestehender Skripte sind eine spätere ausdrückliche
Entscheidung, kein Bestandteil dieser Spezifikationsänderung.

### 18.9 Akzeptanzkriterien für die Implementierung

1. **R1/R2:** Ein frischer Clone enthält alle für Lesen und Wiederherstellung nötigen
   Notes und Metadaten. Nach Löschen sämtlicher Suchcaches lässt sich die Suche aus
   Markdown neu aufbauen; kein Obsidian-Plugin ist erforderlich.
2. **R2/R5:** Suche liefert filterbare, paginierte Treffer; `get_note` liefert den
   vollständigen Inhalt und die Revision. Nach einer Mutation ist der neue Inhalt
   über Suche und Lesezugriff sichtbar.
3. **R3:** Ein beliebiger Reader-Artikel wird durch Archivierung oder Highlighting
   allein nicht importiert. Ein bewusst gewählter, auch bereits archivierter Artikel
   lässt sich als Source mit Original- und Reader-Referenz erfassen.
4. **R3/R4:** Zwei Sources können einen Insight belegen. Der Entwurf trennt Belege,
   Interpretation und Anwendung und wird erst nach bewusster Prüfung aktiv.
5. **R5:** Wiederholung eines Captures mit derselben Reader-ID erzeugt keine zweite
   Source. Ein Retry mit gleicher Request-ID erzeugt keinen zusätzlichen Commit;
   eine abweichende Nutzlast wird abgelehnt.
6. **R5:** Ungültiges Frontmatter, fehlende Link-Ziele und veraltete Revisionen führen
   zu Fehlern ohne Teiländerungen. Zwei konkurrierende Updates überschreiben
   einander nicht. Ein Push-Fehler wird nicht als Remote-Erfolg gemeldet.
7. **R4/R5:** Umbenennen erhält die ID und korrigiert interne Pfadlinks atomar.
   `dry_run` zeigt den geplanten Diff, verändert aber weder Dateien noch Git.
8. **R6:** Falls die CLI weitergeführt wird, akzeptiert und validiert sie das neue
   Schema identisch zum MCP und erhält unbekannte Metadaten bei Updates.
9. **R7:** Eine Probemigration weist für jede ausgewählte Note Herkunft und neue ID
   nach; relevante Inhalte und Links bleiben erhalten. Revert stellt den früheren
   Stand wieder her, der alte Vault bleibt unbeschädigt.
10. **R1/R5:** Ein Token ohne Schreibrecht kann nicht mutieren; Pfade außerhalb des
    Wissensbereichs werden abgelehnt. Credentials gelangen nicht in Notes oder Commits.

Diese Kriterien sind noch nicht abgenommen. Dieser Commit liefert ausschließlich
Dokumentation, keine Implementierung und keine Migration.

### 18.10 Offene Entscheidungen und mögliche spätere Erweiterungen

Vor Implementierung festzulegen sind die endgültigen Repository-Namen, Hosting und
Authentifizierung des MCP, GitHub-API versus Clone als Storage-Adapter, Paketierung
des gemeinsamen Kerns sowie der konkrete Pilotbestand für die Migration.

Mögliche spätere Erweiterungen sind **Ideen, keine zugesagte Version-1-Funktion**:

- Backlink-Ansicht und ähnliche Notes als Hilfe beim Einordnen.
- Semantische Suche/Embeddings als vollständig rekonstruierbarer Zusatzindex.
- Hinweise auf veraltete oder widersprüchliche Aussagen mit überprüfbaren Gründen.
- Unterstütztes Zusammenführen und Aufteilen von Notes mit Erhalt der Referenzen.
- Generierte Topic-Zusammenfassungen und Synthesevorschläge über mehrere Sources.
- Explizit aktivierte Reader-Markierungen als Capture-Auslöser statt manuellem Aufruf.
- Abgeleitete Wissenslandkarten und kontrollierte Retrieval-Workflows.

Die vorhandene Root-`spec.md` wird für diese Änderung bewusst weiterverwendet.
Als separate Dokumentationspflege wird die Migration nach `docs/specs/` mit
SPEC-IDs, YAML-Frontmatter, Index und Backlog gemäß dem globalen Standard vorgeschlagen.
Dabei sind heutiges Verhalten, beschlossene Planung und unverbindliche Ideen sauber
aufzuteilen. Bis dahin wird keine parallele, konkurrierende Spec-Struktur angelegt.
