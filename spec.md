# AI-basiertes Personal Knowledge System

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

Beispielhafte Zielstruktur:

```bash
pkm capture
pkm process inbox
pkm daily create
pkm meeting create
pkm decision create
pkm config show
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

Nicht Teil des MVP, aber vorgesehen:

### `pkm note refine`

Verbessert eine bestehende Notiz AI-gestützt.

### `pkm note link`

Schlägt interne Links und verwandte Notizen vor.

### `pkm note move`

Verschiebt eine Notiz anhand definierter Regeln oder bestätigter Vorschläge.

### `pkm export chat`

Standardisiert Export aus ChatGPT oder anderen Chat-Tools.

### `pkm import session`

Standardisiert Import von Claude-Code- oder Copilot-CLI-Sessions.

### `pkm config show`

Zeigt aktuelle Konfiguration wie Vault-Pfad, Templates-Pfad, Inbox-Pfad.

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
Zentrale Einstiegsnotiz für ein Projekt.

```markdown
---
title: 
type: project
status: active
created: 
updated: 
tags: [project]
area: 
owner: 
related_notes: []
---

# Project

## Goal


## Context


## Current Status


## Key Notes

- 

## Decisions

- 

## Risks / Issues

- 

## Next Steps

- [ ] 

## Links

- 
```

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

## 13. MVP

Der erste lauffähige Stand soll Folgendes umfassen:

* Obsidian Vault mit definierter Ordnerstruktur
* Basis-Templates für die wichtigsten Notiztypen
* einfache Metadaten-Konvention
* Inbox-Workflow
* Export-Workflow für ChatGPT -> Markdown -> Obsidian
* Export-Workflow für Claude-Code- oder Copilot-CLI-Sessions -> Markdown -> Obsidian
* Erstellung von Daily Notes, Meeting Notes und Decision Notes per Claude Code
* 5 bis 8 praxistaugliche Prompts für Claude Code und GitHub Copilot CLI
* optional Git-Repository für Versionierung

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

1. Scope für den MVP finalisieren
2. Ordnerstruktur und Notiztypen bestätigen
3. Templates definieren
4. konkrete Claude-Code-Workflows festlegen
5. Governance-Regeln für Änderungen durch AI festschreiben
6. Beispiel-Vault mit 5 bis 10 Notizen anlegen

