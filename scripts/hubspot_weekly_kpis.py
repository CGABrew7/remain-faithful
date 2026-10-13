#!/usr/bin/env python3
"""
hubspot_weekly_kpis.py

Pulls weekly KPI metrics from HubSpot CRM for the Licensing Pipeline.
Outputs clean JSON with current-week metrics, prior-week metrics, and WoW delta.

Usage:
    HUBSPOT_API_KEY=<token> python hubspot_weekly_kpis.py
    HUBSPOT_API_KEY=<token> python hubspot_weekly_kpis.py --list-owners
    HUBSPOT_API_KEY=<token> python hubspot_weekly_kpis.py --list-stages

Requirements:
    pip install requests
"""

import argparse
import json
import logging
import os
import sys
import time
from datetime import datetime, timedelta
from typing import Optional
from zoneinfo import ZoneInfo

import requests

# ─── Configuration ────────────────────────────────────────────────────────────

HUBSPOT_API_KEY: str = os.environ.get("HUBSPOT_API_KEY", "")
BASE_URL = "https://api.hubapi.com"
PORTAL_ID = "40117269"
PIPELINE_ID = "658429697"
PROPOSE_STAGE_ID = "967626583"
ET = ZoneInfo("America/New_York")

# Populate this dict with HubSpot owner ID → display name to avoid per-deal
# API lookups. Run with --list-owners to generate the entries automatically.
OWNER_MAP: dict[str, str] = {
    # "12345678": "Jane Smith",
    # "87654321": "John Doe",
}

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    datefmt="%H:%M:%S",
    stream=sys.stderr,
)
log = logging.getLogger(__name__)


# ─── Date helpers ─────────────────────────────────────────────────────────────

def week_range_ms(offset: int = 0) -> tuple[int, int]:
    """
    Returns (start_ms, end_ms) in UTC milliseconds for a given week offset in ET.
      offset=0  → current week  Mon 00:00:00 – Sun 23:59:59 ET
      offset=-1 → previous week
    """
    now_et = datetime.now(ET)
    days_since_monday = now_et.weekday()  # Monday=0
    monday = (now_et - timedelta(days=days_since_monday)).replace(
        hour=0, minute=0, second=0, microsecond=0
    )
    monday += timedelta(weeks=offset)
    sunday = monday + timedelta(days=6, hours=23, minutes=59, seconds=59)
    return int(monday.timestamp() * 1000), int(sunday.timestamp() * 1000)


def ms_to_label(ts_ms: int) -> str:
    return datetime.fromtimestamp(ts_ms / 1000, tz=ET).strftime("%Y-%m-%d")


# ─── HTTP helpers ─────────────────────────────────────────────────────────────

def _headers() -> dict[str, str]:
    return {
        "Authorization": f"Bearer {HUBSPOT_API_KEY}",
        "Content-Type": "application/json",
    }


def _get(path: str, params: Optional[dict] = None, retries: int = 6) -> dict:
    url = f"{BASE_URL}{path}"
    backoff = 2
    for attempt in range(retries):
        resp = requests.get(url, headers=_headers(), params=params, timeout=30)
        if resp.status_code == 429:
            wait = int(resp.headers.get("Retry-After", backoff))
            log.warning("Rate limited on GET %s; retrying in %ds (attempt %d/%d)", path, wait, attempt + 1, retries)
            time.sleep(wait)
            backoff = min(backoff * 2, 60)
            continue
        if resp.status_code >= 500:
            log.warning("Server error %d on GET %s; retrying in %ds", resp.status_code, path, backoff)
            time.sleep(backoff)
            backoff = min(backoff * 2, 60)
            continue
        resp.raise_for_status()
        return resp.json()
    raise RuntimeError(f"GET {path} failed after {retries} retries")


def _post(path: str, payload: dict, retries: int = 6) -> dict:
    url = f"{BASE_URL}{path}"
    backoff = 2
    for attempt in range(retries):
        resp = requests.post(url, headers=_headers(), json=payload, timeout=30)
        if resp.status_code == 429:
            wait = int(resp.headers.get("Retry-After", backoff))
            log.warning("Rate limited on POST %s; retrying in %ds (attempt %d/%d)", path, wait, attempt + 1, retries)
            time.sleep(wait)
            backoff = min(backoff * 2, 60)
            continue
        if resp.status_code >= 500:
            log.warning("Server error %d on POST %s; retrying in %ds", resp.status_code, path, backoff)
            time.sleep(backoff)
            backoff = min(backoff * 2, 60)
            continue
        resp.raise_for_status()
        return resp.json()
    raise RuntimeError(f"POST {path} failed after {retries} retries")


# ─── Paginated search ─────────────────────────────────────────────────────────

def search_all(object_type: str, payload: dict) -> list[dict]:
    """Pages through all results from a HubSpot CRM search endpoint."""
    path = f"/crm/v3/objects/{object_type}/search"
    results: list[dict] = []
    after: Optional[str] = None
    page = 0
    while True:
        if after:
            payload["after"] = after
        else:
            payload.pop("after", None)
        data = _post(path, payload)
        batch = data.get("results", [])
        results.extend(batch)
        page += 1
        log.debug("  %s page %d: +%d (total %d)", object_type, page, len(batch), len(results))
        paging = data.get("paging", {})
        after = paging.get("next", {}).get("after")
        if not after:
            break
        time.sleep(0.15)  # stay well under 10 req/s burst limit
    return results


# ─── Pipeline stage resolution ────────────────────────────────────────────────

_stage_cache: dict[str, str] = {}


def get_pipeline_stages() -> dict[str, str]:
    """Returns {stage_id: stage_label} for the Licensing Pipeline."""
    if _stage_cache:
        return _stage_cache
    data = _get(f"/crm/v3/pipelines/deals/{PIPELINE_ID}")
    for stage in data.get("stages", []):
        _stage_cache[str(stage["id"])] = stage["label"]
    log.info("Loaded %d pipeline stages", len(_stage_cache))
    return _stage_cache


def _find_stage_id(keyword: str) -> str:
    stages = get_pipeline_stages()
    for sid, label in stages.items():
        if keyword.lower() in label.lower():
            return sid
    raise ValueError(
        f"Cannot find a stage matching '{keyword}' in pipeline {PIPELINE_ID}.\n"
        f"Available stages: {stages}\n"
        f"Run with --list-stages to inspect."
    )


def closed_won_stage_id() -> str:
    for kw in ("closed won", "won"):
        try:
            return _find_stage_id(kw)
        except ValueError:
            pass
    raise ValueError(f"No Closed Won stage found in pipeline {PIPELINE_ID}. Run --list-stages.")


def closed_lost_stage_id() -> str:
    for kw in ("closed lost", "lost"):
        try:
            return _find_stage_id(kw)
        except ValueError:
            pass
    raise ValueError(f"No Closed Lost stage found in pipeline {PIPELINE_ID}. Run --list-stages.")


# ─── Owner resolution ─────────────────────────────────────────────────────────

_owner_api_cache: dict[str, str] = {}


def resolve_owner(owner_id: str) -> str:
    if not owner_id:
        return "Unassigned"
    if owner_id in OWNER_MAP:
        return OWNER_MAP[owner_id]
    if owner_id in _owner_api_cache:
        return _owner_api_cache[owner_id]
    try:
        data = _get(f"/crm/v3/owners/{owner_id}")
        name = f"{data.get('firstName', '')} {data.get('lastName', '')}".strip()
        _owner_api_cache[owner_id] = name or owner_id
    except Exception as exc:
        log.warning("Could not resolve owner %s: %s", owner_id, exc)
        _owner_api_cache[owner_id] = owner_id
    return _owner_api_cache[owner_id]


def list_all_owners() -> list[dict]:
    """Fetch all CRM owners — use with --list-owners to build OWNER_MAP."""
    data = _get("/crm/v3/owners", params={"limit": 200})
    return data.get("results", [])


# ─── Filter builders ──────────────────────────────────────────────────────────

def _between(prop: str, start_ms: int, end_ms: int) -> dict:
    return {
        "propertyName": prop,
        "operator": "BETWEEN",
        "value": str(start_ms),
        "highValue": str(end_ms),
    }


def _eq(prop: str, value: str) -> dict:
    return {"propertyName": prop, "operator": "EQ", "value": value}


# ─── Metric fetchers ──────────────────────────────────────────────────────────

def get_mqls(start_ms: int, end_ms: int) -> int:
    """Contacts that entered MQL lifecycle stage during the period."""
    results = search_all("contacts", {
        "filterGroups": [{"filters": [
            _between("hs_lifecyclestage_marketingqualifiedlead_date", start_ms, end_ms),
        ]}],
        "properties": ["hs_object_id"],
        "limit": 100,
    })
    return len(results)


def get_sqls(start_ms: int, end_ms: int) -> int:
    """Contacts that entered SQL lifecycle stage during the period."""
    results = search_all("contacts", {
        "filterGroups": [{"filters": [
            _between("hs_lifecyclestage_salesqualifiedlead_date", start_ms, end_ms),
        ]}],
        "properties": ["hs_object_id"],
        "limit": 100,
    })
    return len(results)


def get_appointments(start_ms: int, end_ms: int) -> int:
    """Meetings created/logged in HubSpot during the period."""
    results = search_all("meetings", {
        "filterGroups": [{"filters": [
            _between("createdate", start_ms, end_ms),
        ]}],
        "properties": ["hs_object_id"],
        "limit": 100,
    })
    return len(results)


def _deals_entering_stage(stage_id: str, start_ms: int, end_ms: int) -> list[dict]:
    """
    Deals in the Licensing Pipeline that entered the given stage during the period.
    Relies on HubSpot's auto-generated hs_date_entered_{stage_id} property.
    """
    date_prop = f"hs_date_entered_{stage_id}"
    return search_all("deals", {
        "filterGroups": [{"filters": [
            _eq("pipeline", PIPELINE_ID),
            _between(date_prop, start_ms, end_ms),
        ]}],
        "properties": ["hs_object_id", "dealname", "amount", "hubspot_owner_id", "dealstage"],
        "limit": 100,
    })


def get_agreements_sent(start_ms: int, end_ms: int) -> list[dict]:
    return _deals_entering_stage(PROPOSE_STAGE_ID, start_ms, end_ms)


def get_closed_won_deals(start_ms: int, end_ms: int) -> list[dict]:
    return _deals_entering_stage(closed_won_stage_id(), start_ms, end_ms)


def get_closed_lost_deals(start_ms: int, end_ms: int) -> list[dict]:
    return _deals_entering_stage(closed_lost_stage_id(), start_ms, end_ms)


def get_new_deals(start_ms: int, end_ms: int) -> list[dict]:
    """Deals created in the Licensing Pipeline during the period (for pipeline value)."""
    return search_all("deals", {
        "filterGroups": [{"filters": [
            _eq("pipeline", PIPELINE_ID),
            _between("createdate", start_ms, end_ms),
        ]}],
        "properties": ["hs_object_id", "dealname", "amount", "hubspot_owner_id"],
        "limit": 100,
    })


# ─── Owner breakdown builder ──────────────────────────────────────────────────

def _owner_of(deal: dict) -> str:
    oid = deal.get("properties", {}).get("hubspot_owner_id") or ""
    return resolve_owner(oid)


def _amount_of(deal: dict) -> float:
    try:
        return float(deal.get("properties", {}).get("amount") or 0)
    except (TypeError, ValueError):
        return 0.0


def build_owner_breakdown(
    agreements: list[dict],
    won: list[dict],
    lost: list[dict],
    new_deals: list[dict],
) -> dict[str, dict]:
    owners: dict[str, dict] = {}

    def row(name: str) -> dict:
        if name not in owners:
            owners[name] = {
                "agreements_sent": 0,
                "new_clients_won": 0,
                "closed_lost": 0,
                "pipeline_value_added": 0.0,
            }
        return owners[name]

    for d in agreements:
        row(_owner_of(d))["agreements_sent"] += 1
    for d in won:
        row(_owner_of(d))["new_clients_won"] += 1
    for d in lost:
        row(_owner_of(d))["closed_lost"] += 1
    for d in new_deals:
        row(_owner_of(d))["pipeline_value_added"] += _amount_of(d)

    for v in owners.values():
        v["pipeline_value_added"] = round(v["pipeline_value_added"], 2)

    return dict(sorted(owners.items()))


# ─── Weekly snapshot ──────────────────────────────────────────────────────────

def weekly_snapshot(start_ms: int, end_ms: int) -> dict:
    log.info("  mqls …")
    mqls = get_mqls(start_ms, end_ms)
    log.info("  sqls …")
    sqls = get_sqls(start_ms, end_ms)
    log.info("  appointments …")
    appts = get_appointments(start_ms, end_ms)
    log.info("  agreements sent …")
    agreements = get_agreements_sent(start_ms, end_ms)
    log.info("  closed won …")
    won = get_closed_won_deals(start_ms, end_ms)
    log.info("  closed lost …")
    lost = get_closed_lost_deals(start_ms, end_ms)
    log.info("  new deals …")
    new_deals = get_new_deals(start_ms, end_ms)

    pipeline_value = round(sum(_amount_of(d) for d in new_deals), 2)

    return {
        "mqls_generated": mqls,
        "sqls_generated": sqls,
        "appointments_set": appts,
        "agreements_sent": len(agreements),
        "new_clients_won": len(won),
        "closed_lost": len(lost),
        "pipeline_value_added": pipeline_value,
        "deals_by_owner": build_owner_breakdown(agreements, won, lost, new_deals),
    }


# ─── Main ─────────────────────────────────────────────────────────────────────

SCALAR_KEYS = [
    "mqls_generated", "sqls_generated", "appointments_set",
    "agreements_sent", "new_clients_won", "closed_lost", "pipeline_value_added",
]


def main() -> None:
    parser = argparse.ArgumentParser(
        description="HubSpot weekly KPI report for the Licensing Pipeline"
    )
    parser.add_argument(
        "--list-owners",
        action="store_true",
        help="Print all HubSpot owners as OWNER_MAP entries and exit",
    )
    parser.add_argument(
        "--list-stages",
        action="store_true",
        help="Print all Licensing Pipeline stages and exit",
    )
    args = parser.parse_args()

    if not HUBSPOT_API_KEY:
        sys.exit("ERROR: HUBSPOT_API_KEY environment variable is not set.")

    if args.list_owners:
        print("# Paste into OWNER_MAP in hubspot_weekly_kpis.py")
        print("OWNER_MAP = {")
        for o in list_all_owners():
            name = f"{o.get('firstName', '')} {o.get('lastName', '')}".strip()
            print(f'    "{o["id"]}": "{name}",')
        print("}")
        return

    if args.list_stages:
        print(f"# Stages for pipeline {PIPELINE_ID}")
        for sid, label in get_pipeline_stages().items():
            print(f"  {sid}: {label}")
        return

    curr_start, curr_end = week_range_ms(0)
    prev_start, prev_end = week_range_ms(-1)
    curr_label = ms_to_label(curr_start)
    prev_label = ms_to_label(prev_start)

    log.info("=== Current week (%s) ===", curr_label)
    current = weekly_snapshot(curr_start, curr_end)

    log.info("=== Prior week (%s) ===", prev_label)
    prior = weekly_snapshot(prev_start, prev_end)

    delta = {k: round(current[k] - prior[k], 2) for k in SCALAR_KEYS}

    output = {
        "generated_at": datetime.now(ET).isoformat(),
        "portal_id": PORTAL_ID,
        "pipeline_id": PIPELINE_ID,
        "current_week": {
            "week_start": curr_label,
            "metrics": current,
        },
        "prior_week": {
            "week_start": prev_label,
            "metrics": prior,
        },
        "week_over_week_delta": delta,
    }

    print(json.dumps(output, indent=2))


if __name__ == "__main__":
    main()
