namespace DeviceControlService.Controllers.Dto;

public class CreateSensorRequest
{
    public string Name { get; set; } = string.Empty;
    public string Type { get; set; } = string.Empty;
    public string? Location { get; set; }
}
