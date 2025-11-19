# TigerBeetle Docker Setup

This guide explains how to run TigerBeetle using Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 1.29+

## Quick Start (Single Node)

For development and testing, use the single-node configuration:

```bash
# Start TigerBeetle
docker-compose up -d

# View logs
docker-compose logs -f tigerbeetle

# Stop TigerBeetle
docker-compose down

# Stop and remove all data
docker-compose down -v
```

### Connecting from Go Applications

When TigerBeetle is running via Docker Compose, connect using:

```go
c, err := client.New(client.Config{
    ClusterID:    types.ToUint128(0),
    ReplicaAddrs: []string{"3000"},  // or "localhost:3000"
})
```

## Production Cluster (3 Replicas)

For production-like deployments with high availability, use the cluster configuration:

```bash
# Start 3-node cluster
docker-compose -f docker-compose.cluster.yml up -d

# View logs from all replicas
docker-compose -f docker-compose.cluster.yml logs -f

# Check cluster health
docker-compose -f docker-compose.cluster.yml ps

# Stop cluster
docker-compose -f docker-compose.cluster.yml down

# Stop and remove all data
docker-compose -f docker-compose.cluster.yml down -v
```

### Connecting to Cluster from Go Applications

```go
c, err := client.New(client.Config{
    ClusterID: types.ToUint128(0),
    ReplicaAddrs: []string{
        "localhost:3000",
        "localhost:3001",
        "localhost:3002",
    },
})
```

The client will automatically handle failover if any replica becomes unavailable.

## Configuration Details

### Single Node (`docker-compose.yml`)

- **Port**: 3000
- **Volume**: `tigerbeetle-data` (persistent storage)
- **Use Case**: Development, testing, examples

### Cluster (`docker-compose.cluster.yml`)

- **Replicas**: 3 nodes
- **Ports**: 3000, 3001, 3002
- **Volumes**: Separate volume per replica
- **Use Case**: Production, high availability testing

## Data Persistence

Data is persisted in Docker volumes:

- Single node: `tigerbeetle-data`
- Cluster: `tigerbeetle-data-0`, `tigerbeetle-data-1`, `tigerbeetle-data-2`

### Backup Data

```bash
# Create backup (single node)
docker run --rm -v tigerbeetle-data:/data -v $(pwd):/backup \
  alpine tar czf /backup/tigerbeetle-backup.tar.gz /data

# Restore backup (single node)
docker run --rm -v tigerbeetle-data:/data -v $(pwd):/backup \
  alpine tar xzf /backup/tigerbeetle-backup.tar.gz -C /
```

## Troubleshooting

### Check if TigerBeetle is Running

```bash
# Single node
docker-compose ps

# Cluster
docker-compose -f docker-compose.cluster.yml ps
```

### View Logs

```bash
# Single node - all logs
docker-compose logs

# Single node - follow logs
docker-compose logs -f

# Cluster - specific replica
docker-compose -f docker-compose.cluster.yml logs tigerbeetle-0

# Cluster - all replicas
docker-compose -f docker-compose.cluster.yml logs -f
```

### Connection Refused Errors

If you see connection errors from your Go application:

1. **Check if container is running**:
   ```bash
   docker-compose ps
   ```

2. **Check container logs**:
   ```bash
   docker-compose logs tigerbeetle
   ```

3. **Verify port is accessible**:
   ```bash
   nc -zv localhost 3000
   # or
   telnet localhost 3000
   ```

4. **Check if data directory was initialized**:
   ```bash
   docker-compose exec tigerbeetle ls -la /data
   ```

### Reset Everything

To completely reset and start fresh:

```bash
# Single node
docker-compose down -v
docker-compose up -d

# Cluster
docker-compose -f docker-compose.cluster.yml down -v
docker-compose -f docker-compose.cluster.yml up -d
```

## Running Examples

After starting TigerBeetle with Docker Compose:

```bash
# Run basic transfer example
go run examples/basic_transfer.go

# Run insurance premium collection
go run examples/insurance/premium_collection.go

# Run policy revival example
go run examples/insurance/policy_revival.go
```

## Health Checks

Both configurations include health checks that verify the TigerBeetle server is responsive:

```bash
# Check health status (single node)
docker-compose ps

# Check health status (cluster)
docker-compose -f docker-compose.cluster.yml ps
```

Healthy services will show `(healthy)` in the status.

## Resource Requirements

### Single Node
- **Memory**: 512 MB minimum, 1 GB recommended
- **CPU**: 1 core minimum
- **Disk**: Depends on data volume, starts small and grows

### Cluster (3 replicas)
- **Memory**: 1.5 GB minimum, 3 GB recommended
- **CPU**: 3 cores minimum
- **Disk**: 3x single node requirements

## Production Considerations

For production deployments:

1. **Use the cluster configuration** for high availability
2. **Configure proper volume backups**
3. **Monitor disk usage** as the database grows
4. **Set resource limits** in docker-compose.yml:
   ```yaml
   deploy:
     resources:
       limits:
         cpus: '2'
         memory: 2G
       reservations:
         cpus: '1'
         memory: 512M
   ```

5. **Use Docker Swarm or Kubernetes** for orchestration in production
6. **Configure proper networking** and firewall rules
7. **Enable monitoring and alerting**

## Docker Image Versions

The configurations use `ghcr.io/tigerbeetle/tigerbeetle:latest`.

To pin to a specific version:

```yaml
services:
  tigerbeetle:
    image: ghcr.io/tigerbeetle/tigerbeetle:0.16.65
```

Check available versions at: https://github.com/tigerbeetle/tigerbeetle/pkgs/container/tigerbeetle

## Additional Resources

- [TigerBeetle Documentation](https://docs.tigerbeetle.com/)
- [TigerBeetle GitHub](https://github.com/tigerbeetle/tigerbeetle)
- [Docker Documentation](https://docs.docker.com/)
