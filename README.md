# LRU Memory Cache for Go

This package provides an efficient **Least Recently Used (LRU) Cache** implementation in Go. It supports two modes:
- **LRU with TTL (Time-to-Live)**: Automatically expires items after a specified duration.
- **Simple LRU**: No expiration, only LRU eviction based on cache size.

## Features

- **Thread-safe**: Designed for concurrent access.
- **LRU Eviction**: Ensures the least recently used items are evicted when the cache reaches its size limit.
- **TTL Support** (Optional): Items expire after a configurable time-to-live.

## Installation

```bash
go get github.com/svetoslavdenev/go-mem-cache