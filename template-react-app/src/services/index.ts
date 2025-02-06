/**
 * API Client Library
 *
 * This module provides a robust, isomorphic client for seamless interaction
 * with the API. It handles essential tasks including HTTP request management,
 * response parsing, and authentication workflows.
 *
 * The `ApiClient` class serves as the core interface for API interactions,
 * offering a comprehensive set of methods for data retrieval, manipulation,
 * and analytics processing.
 *
 * The `ApiClientOptions` type defines the configurable parameters for
 * initializing and customizing the API client's behavior.
 */

export type { ApiClientOptions, ErrorResponse, SuccessResponse } from './types/base'

export { default as ApiClient } from './client'
