# GSC / CSC API Reference

> **Canonical source: the GSCode VS Code extension's JSON definitions in [Blakintosh/gscode](https://github.com/Blakintosh/gscode).** Treyarch's local HTML (`docs_modtools/bo3_scriptapifunctions.htm`) is _unreliable_ — it contains errors and omissions. GSCode is community-maintained reverse-engineering with confidence flags per entry; treat it as the live ground truth.

## What's documented (counts)

Pulled directly from the JSON files:

| Language | Function count | Last revision | Revised on (UTC)    | Source file                                                                                                                     |
| -------- | -------------: | ------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| **GSC**  |      **2 191** | rev. 32       | 2026-03-29 12:54:56 | [`server/GSCode.NET/api/t7_api_gsc.json`](https://github.com/Blakintosh/gscode/blob/main/server/GSCode.NET/api/t7_api_gsc.json) |
| **CSC**  |        **801** | rev. 18       | 2026-03-22 22:44:00 | [`server/GSCode.NET/api/t7_api_csc.json`](https://github.com/Blakintosh/gscode/blob/main/server/GSCode.NET/api/t7_api_csc.json) |

(Counts may have grown slightly since this page was written — re-pull the JSON if you want today's numbers.)

## Use it directly via VS Code (recommended)

The friction-free path:

1. Install the [**GSCode** VS Code extension (blakintosh)](https://marketplace.visualstudio.com/items?itemName=blakintosh.gscode).
2. Open any `.gsc` / `.csc` file.
3. Hover any built-in function for its description, parameters, and example. IntelliSense autocompletes by name. **`F12` (go-to-definition) only works on user / script functions** — for built-in API entries, hover is the only path; there's no source-line to jump to.

That's the answer for 90 % of "which function does X" questions. The rest of this page is for when you need to _browse_ / _grep_ / _cite_ the API outside VS Code.

## How an entry is structured

Each function is a JSON object with this shape (real example: **`SetLightingState`**):

```json
{
  "name": "SetLightingState",
  "description": "Changes lighting state for the map.",
  "overloads": [
    {
      "calledOn": null,
      "parameters": [
        {
          "name": "newLightState",
          "description": "new state to change to. Lighting state defaults to 1 at start of game.",
          "mandatory": true,
          "type": { "dataType": "int", "isArray": false }
        }
      ]
    }
  ],
  "flags": ["processed"],
  "example": "SetLightingState( 2 )",
  "confidence": "high"
}
```

Field-by-field:

| Key           | Meaning                                                                                                                                                 |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`        | The function name as you call it from script (case-sensitive in the JSON; GSC itself is case-insensitive at the language level).                        |
| `description` | Short human description.                                                                                                                                |
| `overloads`   | Array — most functions have **one** overload, but some have multiple parameter shapes. Each overload has its own `calledOn` / `parameters` / `returns`. |
| `calledOn`    | If `null` → free function (`Foo(x)`). If an object describing a type → method-style call on that type (`self Foo(x)` where `self` is e.g. a `player`).  |
| `parameters`  | Array of `{name, description, mandatory, type}`. Type holds `dataType` + `isArray` flag.                                                                |
| `returns`     | Object with `name`, `description`, `type`, `void`. Optional — sometimes elided.                                                                         |
| `flags`       | Tags — see flag table below.                                                                                                                            |
| `example`     | One-line usage example.                                                                                                                                 |
| `confidence`  | `high` / `medium` / `low` — the GSCode maintainers' confidence in this entry's accuracy.                                                                |

### Flags

From the GSC JSON (counts):

| Flag        | Count | What it means                                                    |
| ----------- | ----: | ---------------------------------------------------------------- |
| `processed` | 2 033 | Default state — entry has been processed by the GSCode pipeline. |
| `verified`  |   156 | Manually verified accurate.                                      |
| `unlisted`  |     3 | Hidden from IntelliSense (likely deprecated / engine-internal).  |
| `Broken`    |     1 | Known broken in this engine version — don't use.                 |

**Most useful**: when an entry is `verified`, trust it; when it's only `processed` and `confidence: low`, treat the description as a hint rather than gospel — check `share/raw/scripts/` for actual usage to confirm.

### Confidence breakdown (GSC, today)

| Confidence | Count |
| ---------- | ----: |
| `high`     | 1 291 |
| `medium`   |   685 |
| `low`      |    80 |
| (unset)    |   135 |

So roughly 60 % of the API is high-confidence; 30 % medium; the rest is shakier. Cross-check anything labelled low / unset against actual Treyarch usage in `share/raw/scripts/`.

## Free functions vs method-style functions

`calledOn: null` means it's called as a free function: `Foo(args)`.

Anything else means it's called as a **method on a target type** (the `self` of the call). The major target types in the GSC API:

| `calledOn` type                  | Function count | Typical use                                                   |
| -------------------------------- | -------------: | ------------------------------------------------------------- |
| _(free)_                         |            862 | Math, utility, world queries.                                 |
| `player`                         |            418 | Score, weapons, perks, HUD, inventory.                        |
| `entity`                         |      254 + 238 | Generic entity ops (origin/angles, hide/show, link, animate). |
| `vehicle`                        |            127 | Vehicle-specific ops.                                         |
| `actor`                          |            122 | AI behaviour ops.                                             |
| `sentient`                       |             28 | Damage/health on anything that can take damage.               |
| `turret`                         |             24 | Turret-specific.                                              |
| `hud_element`                    |             20 | HUD element manipulation.                                     |
| `client`                         |             15 | Client connection state.                                      |
| `trigger`                        |             12 | Trigger-specific (entity is the trigger).                     |
| `script_model` / `script_origin` |             11 | Script-placed model / origin entities.                        |
| `light`                          |             10 | Light entities.                                               |
| `ai_or_player`                   |              8 | Functions valid on either AI or player.                       |
| `ai`                             |              5 | AI-specific (narrower than `actor`).                          |

You write `calledOn` calls as `<target> <Function>(args)`:

```c
self GiveWeapon( "iw8_ak47_zm" );      // self == player
trig SetCursorHint( "HINT_ACTIVATE" ); // trig == trigger entity
veh SetSpeedImmediate( 100 );          // veh == vehicle
```

When in doubt, hover in GSCode — it shows the `calledOn` type in the hover popover.

## How to use the JSON outside VS Code

### Pull both files

```bash
curl -O https://raw.githubusercontent.com/Blakintosh/gscode/main/server/GSCode.NET/api/t7_api_gsc.json
curl -O https://raw.githubusercontent.com/Blakintosh/gscode/main/server/GSCode.NET/api/t7_api_csc.json
```

### Search the JSON with `jq`

```bash
# All function names containing "FX"
jq -r '.api[] | select(.name | test("FX")) | .name' t7_api_gsc.json

# Show one entry in full
jq '.api[] | select(.name == "SetLightingState")' t7_api_gsc.json

# Count entries by calledOn type
jq -r '.api[].overloads[].calledOn // "<free>" | if type == "object" then .name else . end' t7_api_gsc.json | sort | uniq -c | sort -rn
```

### Search inline (no jq)

```bash
grep -A 10 '"name": "SetLightingState"' t7_api_gsc.json
```

## When the JSON is wrong or thin

GSCode is a community project — entries can be missing, mislabelled, or wrong. Cross-references when something looks off:

- **`share/raw/scripts/`** — Treyarch's actual GSC source ships in your install. `grep -rn "FunctionName(" share/raw/scripts/` to see how they use it.
- **[`shiversoftdev/t7-source`](https://github.com/shiversoftdev/t7-source)** — broader Treyarch source / engine dump. Closest to engine ground truth.
- **The HTML / web API references** ([`docs_modtools/bo3_scriptapifunctions.htm`](file:///D:/SteamLibrary/steamapps/common/Call%20of%20Duty%20Black%20Ops%20III/docs_modtools/bo3_scriptapifunctions.htm) locally and [BO3 Source Code Explorer (zeroy)](https://bo3explorer.zeroy.com/) on the web) are essentially the same content surfaced in two ways. **Both are known to have errors and omissions** — use only as a last resort, prefer the GSCode JSON.
- **Community Discords** — for "is this function actually broken in version Y" questions.

### File a fix upstream

If you confirm the JSON has an error, file an issue or PR against [Blakintosh/gscode](https://github.com/Blakintosh/gscode/issues) — that's how the canonical truth gets updated for everyone using GSCode.

## Direct links — keep these handy

- **GSC JSON (raw)**: <https://raw.githubusercontent.com/Blakintosh/gscode/main/server/GSCode.NET/api/t7_api_gsc.json>
- **CSC JSON (raw)**: <https://raw.githubusercontent.com/Blakintosh/gscode/main/server/GSCode.NET/api/t7_api_csc.json>
- **GSC JSON (browsable on GitHub)**: <https://github.com/Blakintosh/gscode/blob/main/server/GSCode.NET/api/t7_api_gsc.json>
- **CSC JSON (browsable)**: <https://github.com/Blakintosh/gscode/blob/main/server/GSCode.NET/api/t7_api_csc.json>
- **GSCode VS Code extension**: <https://marketplace.visualstudio.com/items?itemName=blakintosh.gscode>
- **GSCode repo**: <https://github.com/Blakintosh/gscode>
- **File a fix**: <https://github.com/Blakintosh/gscode/issues>

## Adjacent references

- The non-canonical local HTML: [`docs_modtools/bo3_scriptapifunctions.htm`](file:///D:/SteamLibrary/steamapps/common/Call%20of%20Duty%20Black%20Ops%20III/docs_modtools/bo3_scriptapifunctions.htm) — useful only when GSCode is unreachable; cross-check anything you read here.
- docs_modtools/GSC_Language.pdf — language _reference_ (syntax, threads, namespaces) rather than function API. Source for the planned `gsc-language.md` page.
- `share/raw/scripts/` — the in-repo mirror of shipped script source. Grep for usage examples in real Treyarch context.
