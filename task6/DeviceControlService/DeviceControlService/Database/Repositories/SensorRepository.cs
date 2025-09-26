using DeviceControlService.Database.Models;
using Microsoft.EntityFrameworkCore;
using System;

namespace DeviceControlService.Database.Repositories;

public class SensorRepository : ISensorRepository
{
    private readonly DeviceContext _context;

    public SensorRepository(DeviceContext context)
    {
        _context = context;
    }

    public async Task<Sensor?> GetByIdAsync(Guid id)
    {
        return await _context.Sensors.FindAsync(id);
    }

    public async Task<IEnumerable<Sensor>> GetAllAsync()
    {
        return await _context.Sensors
            .OrderBy(s => s.Name)
            .ToListAsync();
    }

    public async Task<Sensor> CreateAsync(Sensor sensor)
    {
        _context.Sensors.Add(sensor);
        await _context.SaveChangesAsync();
        return sensor;
    }

    public async Task<Sensor?> UpdateAsync(Sensor sensor)
    {
        var existingSensor = await _context.Sensors.FindAsync(sensor.Id);
        if (existingSensor == null)
            return null;

        existingSensor.Name = sensor.Name;
        existingSensor.Type = sensor.Type;
        existingSensor.Location = sensor.Location;
        existingSensor.IsActive = sensor.IsActive;
        existingSensor.UpdatedAt = DateTime.UtcNow;

        await _context.SaveChangesAsync();
        return existingSensor;
    }

    public async Task<bool> DeleteAsync(Guid id)
    {
        var sensor = await _context.Sensors.FindAsync(id);
        if (sensor == null)
            return false;

        _context.Sensors.Remove(sensor);
        await _context.SaveChangesAsync();
        return true;
    }

    public async Task<bool> ExistsByNameAsync(string name)
    {
        return await _context.Sensors
            .AnyAsync(s => s.Name.ToLower() == name.ToLower());
    }

    public async Task<IEnumerable<Sensor>> GetAllSensorsAsync()
    {
        return await _context.Sensors
            .Where(s => s.IsActive)
            .OrderBy(s => s.Name)
            .ToListAsync();
    }
}
