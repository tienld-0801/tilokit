package csharp

const (
	AspNetCsproj = `<Project Sdk="Microsoft.NET.Sdk.Web">

  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <Nullable>enable</Nullable>
    <ImplicitUsings>enable</ImplicitUsings>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="Microsoft.AspNetCore.OpenApi" Version="8.0.0" />
    <PackageReference Include="Swashbuckle.AspNetCore" Version="6.4.0" />
  </ItemGroup>

</Project>`

	AspNetProgram = `var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();

app.MapGet("/", () => new { message = "Hello World from {{.ProjectName}}!" });
app.MapGet("/health", () => new { status = "ok" });

app.Run();`

	AspNetController = `using Microsoft.AspNetCore.Mvc;

namespace {{.ProjectName}}.Controllers;

[ApiController]
[Route("[controller]")]
public class HomeController : ControllerBase
{
    [HttpGet]
    public IActionResult Get()
    {
        return Ok(new { message = "Hello World from {{.ProjectName}}!" });
    }

    [HttpGet("health")]
    public IActionResult Health()
    {
        return Ok(new { status = "ok" });
    }
}`
)
