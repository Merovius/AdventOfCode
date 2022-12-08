#include <fstream>
#include <iostream>
#include <string>
#include <tuple>
#include <unordered_set>
#include <vector>

// Using abseil for error-handling and string utilities. Compile with
// 	g++ -std=c++20 -o day13 day13.cpp -l absl_status -l absl_statusor -l absl_strings
#include <absl/status/statusor.h>
#include <absl/strings/numbers.h>
#include <absl/strings/str_split.h>
#include <absl/strings/string_view.h>

enum Axis {
  AxisX,
  AxisY,
};

struct Fold {
  Axis on;
  int value;
};

typedef std::pair<int, int> Point;

struct PointHash {
  std::size_t operator()(const Point &p) const {
    return std::hash<int>()(p.first) ^ std::hash<int>()(p.second);
  }
};

typedef std::unordered_set<Point, PointHash> PointSet;

struct Data {
  PointSet points;
  std::vector<Fold> folds;
};

absl::StatusOr<Data> read(const std::string &filename) {
  Data data;

  std::ifstream file;
  try {
    file.open(filename);
    std::string line;
    bool found_sep = false;
    while (std::getline(file, line)) {
      if (line.empty()) {
        found_sep = true;
        break;
      }
      std::vector<absl::string_view> v = absl::StrSplit(line, ",");
      if (v.size() != 2) {
        return absl::InvalidArgumentError("file contains invalid point");
      }
      int x, y;
      if (!absl::SimpleAtoi(v[0], &x)) {
        return absl::InvalidArgumentError("file contains invalid point");
      }
      if (!absl::SimpleAtoi(v[1], &y)) {
        return absl::InvalidArgumentError("file contains invalid point");
      }
      data.points.insert(std::make_pair(x, y));
    }
    if (!found_sep) {
      return absl::InvalidArgumentError("file contains no empty line");
    }
    while (std::getline(file, line)) {
      auto l = absl::string_view(line);
      if (!absl::ConsumePrefix(&l, "fold along ")) {
        return absl::InvalidArgumentError("fold instruction without prefix");
      }
      Fold f;
      std::vector<absl::string_view> v = absl::StrSplit(l, "=");
      if (v.size() != 2) {
        return absl::InvalidArgumentError(
            "fold instruction without assignment");
      }
      if (v[0] == "x") {
        f.on = AxisX;
      } else if (v[0] == "y") {
        f.on = AxisY;
      } else {
        return absl::InvalidArgumentError("fold instruction with invalid axis");
      }
      if (!absl::SimpleAtoi(v[1], &f.value)) {
        return absl::InvalidArgumentError(
            "fold instruction with invalid value");
      }
      data.folds.push_back(f);
    }
  } catch (...) {
    return absl::InternalError("could not read file");
  };
  return data;
}

Point apply_to_point(Point p, Fold f) {
  switch (f.on) {
  case AxisX:
    return std::make_pair(p.first < f.value ? p.first : 2 * f.value - p.first,
                          p.second);
  case AxisY:
    return std::make_pair(p.first, p.second < f.value ? p.second
                                                      : 2 * f.value - p.second);
  }
  throw std::runtime_error("invalid folding axis");
}

void apply(PointSet *points, Fold f) {
  std::vector<Point> next;
  for (auto p : *points) {
    next.push_back(apply_to_point(p, f));
  }
  points->clear();
  for (auto p : next) {
    points->insert(p);
  }
}

void print_paper(const PointSet &points) {
  int mx = 0;
  int my = 0;
  for (auto p : points) {
    mx = p.first > mx ? p.first : mx;
    my = p.second > my ? p.second : my;
  }
  Point p;
  for (p.second = 0; p.second <= my; p.second++) {
    for (p.first = 0; p.first <= mx; p.first++) {
      if (points.contains(p)) {
        std::cout << "•";
      } else {
        std::cout << " ";
      }
    }
    std::cout << std::endl;
  }
}

int main() {
  auto result = read("input.txt");
  if (!result.ok()) {
    std::cerr << result.status() << std::endl;
    return 1;
  }
  auto data = *result;

  apply(&data.points, data.folds.at(0));
  std::cout << "After first fold, there are " << data.points.size() << " points"
            << std::endl;
  for (int i = 1; i < data.folds.size(); i++) {
    apply(&data.points, data.folds.at(i));
  }
  std::cout << "After all folds, the paper looks like:" << std::endl;
  print_paper(data.points);
}
