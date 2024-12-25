const std = @import("std");

const Data = struct { task_id: []const u8, client_thumbprint: []const u8, search_item_timestamp: []const u8, search_item: []const u8, search_results: []const u8, search_results_captured: []const u8, sum_bundle: []const u8, num_outlier: []const u8, percent_complete: []const u8, p0: []const u8, p25: []const u8, p50: []const u8, p75: []const u8, p100: []const u8, mean: []const u8, std: []const u8 };
pub fn main() !void {
    var arena = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();
    const allocator = arena.allocator();
    var client = std.http.Client{ .allocator = allocator };
    defer client.deinit();

    const endpoint = "https://uselessfacts.jsph.pl/api/v2/facts/random";
    const uri = try std.Uri.parse(endpoint);

    const headers = std.http.Client.Request.Headers{
        .content_type = std.http.Client.Request.Headers.Value{
            .override = "application/json",
        },
    };

    const server_header_buffer: []u8 = try allocator.alloc(u8, 8 * 1024 * 4);

    var req = try client.open(.GET, uri, std.http.Client.RequestOptions{
        .server_header_buffer = server_header_buffer,
        .headers = headers,
    });
    defer req.deinit();

    try req.send();
    try req.wait();

    const json_str = try req.reader().readAllAlloc(allocator, std.math.maxInt(usize));
    defer allocator.free(json_str);

    const data = try std.json.parseFromSlice(Data, allocator, json_str, .{});
    defer data.deinit();

    std.debug.print("fact! {any}\n", .{data.value});
}
