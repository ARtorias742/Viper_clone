
### Explanation of Extensions

1. **JSON Support**:
   - Added `unmarshalJSON` to parse JSON data in `ReadInConfig`.
   - Updated `WriteConfig` to handle JSON encoding with proper indentation.

2. **Command-Line Flags**:
   - Integrated `pflag` by adding a `flags` field to the `Viper` struct.
   - Added `BindPFlags` to associate a `pflag.FlagSet` with Viper.
   - Modified `Get` to prioritize flag values over config file values.

3. **Environment Variable Binding**:
   - Added `AutomaticEnv` (currently a no-op, as binding is handled in `Get`).
   - Updated `Get` to check environment variables after config file values but before defaults.
   - Environment variables are normalized (e.g., `port.number` becomes `PORT_NUMBER`).

### How to Test the Extensions

1. **Set Up**:
   - Replace the code in each file with the updated versions above.
   - Run `go mod tidy` if you haven’t already to ensure dependencies are resolved.

2. **Run the Example**:
   - Test JSON support:
     ```bash
     go run cmd/main.go