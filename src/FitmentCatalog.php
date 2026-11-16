<?php

declare(strict_types=1);

namespace CustomGraphics\DirtBikeFitment;

use RuntimeException;

final class FitmentCatalog
{
    /** @var array{updated: string, bikes: list<array{id: string, brand: string, model: string, reference_year: int, category: string, fits: list<string>}>}|null */
    private static ?array $cache = null;

    /**
     * @return array{updated: string, bikes: list<array{id: string, brand: string, model: string, reference_year: int, category: string, fits: list<string>}>}
     */
    public static function raw(): array
    {
        if (self::$cache !== null) {
            return self::$cache;
        }
        $candidates = [
            dirname(__DIR__) . '/fitment.json',
            dirname(__DIR__) . '/data/fitment.json',
        ];
        foreach ($candidates as $path) {
            if (is_file($path)) {
                $json = file_get_contents($path);
                if ($json === false) {
                    throw new RuntimeException('Unable to read fitment.json');
                }
                /** @var array{updated: string, bikes: list<array{id: string, brand: string, model: string, reference_year: int, category: string, fits: list<string>}>} $decoded */
                $decoded = json_decode($json, true, 512, JSON_THROW_ON_ERROR);
                self::$cache = $decoded;
                return $decoded;
            }
        }
        throw new RuntimeException('fitment.json dataset not found');
    }

    /**
     * @return list<array{id: string, brand: string, model: string, reference_year: int, category: string, fits: list<string>}>
     */
    public static function bikes(): array
    {
        return self::raw()['bikes'];
    }

    public static function updated(): string
    {
        return self::raw()['updated'];
    }

    /**
     * Find the reference bike whose plastics and graphics kit fit the given model and optional year.
     *
     * @return array{id: string, brand: string, model: string, reference_year: int, category: string, fits: list<string>}|null
     */
    public static function findKit(string $model, int|string|null $year = null): ?array
    {
        $normModel = self::normalize($model);
        if ($normModel === '') {
            return null;
        }
        $yearStr = $year === null ? null : (string) $year;
        foreach (self::bikes() as $bike) {
            foreach ($bike['fits'] as $fit) {
                $spacePos = strpos($fit, ' ');
                if ($spacePos === false) {
                    continue;
                }
                $fitYear = substr($fit, 0, $spacePos);
                $fitModel = substr($fit, $spacePos + 1);
                if (self::normalize($fitModel) === $normModel && ($yearStr === null || $fitYear === $yearStr)) {
                    return $bike;
                }
            }
        }
        return null;
    }

    /**
     * Return every model year that shares plastics with the queried bike.
     *
     * @return list<string>
     */
    public static function sharedFits(string $model, int|string|null $year = null): array
    {
        $bike = self::findKit($model, $year);
        return $bike !== null ? $bike['fits'] : [];
    }

    /**
     * True when one graphics kit template fits both bikes.
     */
    public static function sameKit(string $modelA, int|string|null $yearA, string $modelB, int|string|null $yearB): bool
    {
        $kitA = self::findKit($modelA, $yearA);
        if ($kitA === null) {
            return false;
        }
        $kitB = self::findKit($modelB, $yearB);
        return $kitB !== null && $kitA['id'] === $kitB['id'];
    }

    private static function normalize(string $value): string
    {
        return (string) preg_replace('/[^A-Z0-9]/', '', strtoupper($value));
    }
}
