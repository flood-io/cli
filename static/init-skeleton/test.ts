import { step, TestSettings, Until, By, MouseButtons, Device, Driver } from '@flood/chrome'
import * as assert from 'assert'
export const settings: TestSettings = {
	loopCount: -1,
	device: Device.iPadLandscape,
	userAgent: 'flood-chrome-test',
	// clearCache: true,
	disableCache: true,
	actionDelay: 0.5,
	stepDelay: 2.5,
}

/**
 * {{.Name}}
 * Version: 1.0
 */
export default () => {
	step('Test: Start', async (browser: Driver) => {
		await browser.visit('{{.URL}}')

		assert(true, 'congratulations!')
	})
}
