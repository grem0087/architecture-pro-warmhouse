using DeviceControlService.Controllers.Dto;
using DeviceControlService.Database.Models;
using DeviceControlService.Database.Repositories;
using Microsoft.AspNetCore.Mvc;

namespace DeviceControlService.Controllers;

[ApiController]
[Route("api/[controller]")]
public class SensorsController : ControllerBase
{
    private readonly ISensorRepository _sensorRepository;

    public SensorsController(ISensorRepository sensorRepository)
    {
        _sensorRepository = sensorRepository;
    }

    [HttpPost]
    public async Task<ActionResult> CreateSensor([FromBody] CreateSensorRequest request)
    {
        var sensor = new Sensor
        {
            Name = request.Name,
            Type = request.Type,
            Location = request.Location,
            IsActive = true,
            CreatedAt = DateTime.UtcNow
        };

        var createdSensor = await _sensorRepository.CreateAsync(sensor);

        return Ok(new { id = sensor.Id });
    }

    [HttpGet]
    public async Task<ActionResult<IEnumerable<Sensor>>> GetAllSensors()
    {
        var sensors = await _sensorRepository.GetAllSensorsAsync();
        return Ok(sensors);
    }
}
