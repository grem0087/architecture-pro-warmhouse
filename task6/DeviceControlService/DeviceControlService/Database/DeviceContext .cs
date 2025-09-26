using DeviceControlService.Database.Models;
using Microsoft.EntityFrameworkCore;

namespace DeviceControlService.Database;
public class DeviceContext : DbContext
{
    public DeviceContext(DbContextOptions<DeviceContext> options) : base(options)
    {
    }

    public DbSet<Sensor> Sensors { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Sensor>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.Property(e => e.Name).IsRequired().HasMaxLength(100);
            entity.Property(e => e.Type).IsRequired().HasMaxLength(50);
            entity.Property(e => e.Location).HasMaxLength(200);
            entity.Property(e => e.CreatedAt).HasDefaultValueSql("NOW()");

            // Индекс для быстрого поиска по имени
            entity.HasIndex(e => e.Name).IsUnique();

            // Индекс для фильтрации по активности
            entity.HasIndex(e => e.IsActive);
        });
    }
}

