import { useState, useRef } from "react";
import {
  HStack,
  VStack,
  Icon,
  Text,
  Spinner,
  Box,
  Stack,
  Button,
  Input,
  Heading,
  Flex,
} from "@chakra-ui/react";
import { FiPlusSquare, FiInfo } from "react-icons/fi";
import { FaArrowLeft } from "react-icons/fa6";
import {
  fetchExternalSystems,
  deleteExternalSystem,
  addExternalSystem,
} from "@/api/externalsystems";
import { useQuery, useQueryClient, useMutation } from "@tanstack/react-query";
import { toaster, Toaster } from "@/components/ui/toaster";
import { Tooltip } from "@/components/ui/tooltip";
import { Link } from "react-router-dom";

export default function ExternalSystemsPage() {
  const [loadingIndex, setLoadingIndex] = useState<number | null>(null);
  const [externalName, setExternalName] = useState("");
  const { data, error, isLoading } = useQuery({
    queryKey: ["externalSystems"],
    queryFn: fetchExternalSystems,
  });

  const queryClient = useQueryClient();
  const ref = useRef<HTMLInputElement>(null);
  const deleteMutation = useMutation({
    mutationFn: (systemName: string) => deleteExternalSystem(systemName),
    onSuccess: () => {
      toaster.create({
        title: "Successful deletion!",
        description: "The system you have selected has been deleted",
        type: "success",
        duration: 2000,
      });

      queryClient.invalidateQueries({ queryKey: ["externalSystems"] });
    },
    onError: () => {
      toaster.create({
        title: "Failed to delete!",
        description: "The system you have selected has not been deleted",
        type: "error",
        duration: 2000,
      });
    },
    onSettled: () => {
      setLoadingIndex(null);
    },
  });

  const addMutation = useMutation({
    mutationFn: (systemName: string) => addExternalSystem(systemName),
    onSuccess: () => {
      toaster.create({
        title: "System added!",
        description: "The external system was successfully added.",
        type: "success",
        duration: 2000,
      });
      setExternalName("");
      queryClient.invalidateQueries({ queryKey: ["externalSystems"] });
    },
    onError: () => {
      toaster.create({
        title: "Failed to add system!",
        description: "An error occurred while adding the external system.",
        type: "error",
        duration: 2000,
      });
    },
  });

  const handleDelete = (name: string, index: number) => {
    setLoadingIndex(index);
    deleteMutation.mutate(name);
  };
  return (
    <Flex
      minH="100vh"
      align="center"
      justify="center"
      bg="gray.50"
      px={4}
      py={10}
    >
      <Box
        bgGradient="linear(to-b, white, gray.50)"
        px={8}
        py={4}
        rounded="lg"
        boxShadow="lg"
        w="70%"
        textAlign="center"
      >
        <VStack gap={6} align="stretch">
          <Box mt={4} textAlign="center">
            <Link to="/dashboard">
              <HStack
                justify="center"
                gap={1}
                color="gray.600"
                _hover={{ color: "blue.500" }}
              >
                <Icon as={FaArrowLeft} />
                <Text fontSize="sm">Back to dashboard</Text>
              </HStack>
            </Link>
          </Box>
          <HStack gap={2} align="center">
            <Icon as={FiPlusSquare} boxSize={6} />
            <Heading size="md">External Systems</Heading>
          </HStack>

          <Box>
            <HStack>
              <Input
                ref={ref}
                value={externalName}
                onChange={(e) => setExternalName(e.target.value)}
                placeholder="Enter system name"
                onKeyDown={(e) => {
                  if (e.key === "Enter" && externalName.trim()) {
                    addMutation.mutate(externalName.trim());
                  }
                }}
              />
              <Button
                onClick={() => addMutation.mutate(externalName.trim())}
                disabled={!externalName.trim()}
              >
                Add New System
              </Button>
            </HStack>

            {addMutation.data?.access_key && (
              <Box
                mt={4}
                p={4}
                borderWidth="1px"
                borderRadius="md"
                textAlign="center"
              >
                <Text fontWeight="bold">New Access Token:</Text>
                <HStack
                  w="100%"
                  justify="center"
                  align="center"
                  mt={2}
                >
                  <Text color="green.600">{addMutation.data.access_key}</Text>
                  <Tooltip
                    content="This is your API key. Make sure to copy it now — you won’t be able to view it again!"
                    aria-label="Copy your API key reminder"
                    positioning={{ placement: "right-end" }}
                    showArrow
                    openDelay={0}
                    closeDelay={0}
                    open={true}
                  >
                    <Icon as={FiInfo}></Icon>
                  </Tooltip>
                </HStack>
              </Box>
            )}
          </Box>

          {isLoading && (
            <Box textAlign="center" py={4}>
              <Spinner />
              <Text mt={2}>Loading…</Text>
            </Box>
          )}
          {error && (
            <Box color="red.500" py={4}>
              Failed to load systems.
            </Box>
          )}

          {!isLoading && !error && (
            <Stack gap={3}>
              {data?.slice().map((name: string, index: number) => (
                <Box
                  key={index}
                  p={3}
                  borderWidth="1px"
                  borderRadius="lg"
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Text fontWeight="medium">{name}</Text>
                  <Button
                    size="sm"
                    colorScheme="red"
                    onClick={() => handleDelete(name, index)}
                    loading={loadingIndex === index}
                  >
                    Delete
                  </Button>
                </Box>
              ))}
            </Stack>
          )}

          <Toaster />
        </VStack>
      </Box>
    </Flex>
  );
}
