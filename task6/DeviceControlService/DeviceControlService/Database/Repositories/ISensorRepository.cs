using DeviceControlService.Controllers.Dto;
using DeviceControlService.Database.Models;

namespace DeviceControlService.Database.Repositories;

public interface ISensorRepository
{
    Task<Sensor?> GetByIdAsync(Guid id);
    Task<IEnumerable<Sensor>> GetAllAsync();
    Task<Sensor> CreateAsync(Sensor sensor);
    Task<Sensor?> UpdateAsync(Sensor sensor);
    Task<bool> DeleteAsync(Guid id);
    Task<bool> ExistsByNameAsync(string name);
    Task<IEnumerable<Sensor>> GetAllSensorsAsync();
}
